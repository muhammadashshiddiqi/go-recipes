package repository

import (
	"context"
	"errors"
	"fmt"
	"go-recipe/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const recipeColumns = `id, user_id, title, description, ingredients, steps, cook_time_minutes, servings, created_at, updated_at`

type RecipeRepository struct{ db *pgxpool.Pool }

func NewRecipeRepository(db *pgxpool.Pool) *RecipeRepository {
	return &RecipeRepository{db: db}
}

func scanRecipe(row pgx.Row) (*model.Recipe, error) {
	var r model.Recipe
	err := row.Scan(&r.ID, &r.UserID, &r.Title, &r.Description, &r.Ingredients, &r.Steps, &r.CookTimeMinutes, &r.Servings, &r.CreatedAt, &r.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &r, nil
}

func (r *RecipeRepository) Create(ctx context.Context, userID int64, in model.RecipeInput) (*model.Recipe, error) {
	return scanRecipe(r.db.QueryRow(ctx, `
	INSERT INTO recipes (user_id, title, description, ingredients, steps, cook_time_minutes, servings)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING
	`+recipeColumns,
		userID, in.Title, in.Description, in.Ingredients, in.Steps, in.CookTimeMinutes, in.Servings))

}

func (r *RecipeRepository) GetByID(ctx context.Context, id int64) (*model.Recipe, error) {
	return scanRecipe(r.db.QueryRow(ctx, `SELECT `+recipeColumns+` FROM recipes WHERE id = $1`, id))
}

func (r *RecipeRepository) List(ctx context.Context, search string, limit, offset int) ([]model.Recipe, int, error) {
	where, args := "", []any{}
	if search != "" {
		where = "WHERE title ILIKE $1"
		args = append(args, "%"+search+"%")
	}

	var total int
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM recipes "+where, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`SELECT %s FROM recipes %s ORDER BY created_at DESC, id DESC LIMIT $%d OFFSET $%d`, recipeColumns, where, len(args)+1, len(args)+2)
	rows, err := r.db.Query(ctx, query, append(args, limit, offset)...)

	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	recipes := make([]model.Recipe, 0, limit)
	for rows.Next() {
		rec, err := scanRecipe(rows)
		if err != nil {
			return nil, 0, err
		}
		recipes = append(recipes, *rec)
	}
	return recipes, total, rows.Err()
}

func (r *RecipeRepository) Update(ctx context.Context, id, userID int64, in model.RecipeInput) (*model.Recipe, error) {
	query := `UPDATE recipes SET 
		title=$1, description=$2, ingredients=$3, steps=$4, 
		cook_time_minutes=$5, servings=$6, updated_at=NOW()
		WHERE id=$7 AND user_id=$8 RETURNING ` + recipeColumns

	return scanRecipe(r.db.QueryRow(ctx, query, in.Title, in.Description, in.Ingredients, in.Steps, in.CookTimeMinutes, in.Servings, id, userID))
}

func (r *RecipeRepository) Delete(ctx context.Context, id, userID int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM recipes WHERE id=$1 AND user_id=$2`, id, userID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

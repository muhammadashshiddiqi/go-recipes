package handler

import (
	"errors"
	"go-recipe/internal/auth"
	"go-recipe/internal/model"
	"go-recipe/internal/repository"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type RecipeHandler struct{ repo *repository.RecipeRepository }

func NewRecipeHandler(repo *repository.RecipeRepository) *RecipeHandler {
	return &RecipeHandler{repo: repo}
}

func parseID(c *fiber.Ctx) (int64, error) {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)

	if err != nil || id <= 0 {
		return 0, fiber.NewError(fiber.StatusBadRequest, "id not valid")
	}
	return id, nil
}

func parseInput(c *fiber.Ctx) (model.RecipeInput, error) {
	var in model.RecipeInput
	if err := c.BodyParser(&in); err != nil {
		return in, fiber.NewError(fiber.StatusBadRequest, "body request not valid")
	}
	if msg := in.Validate(); msg != "" {
		return in, fiber.NewError(fiber.StatusUnprocessableEntity, msg)
	}
	return in, nil
}

func (h *RecipeHandler) Create(c *fiber.Ctx) error {
	in, err := parseInput(c)
	if err != nil {
		return err
	}

	recipe, err := h.repo.Create(c.Context(), auth.UserID(c), in)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": recipe})
}

func (h *RecipeHandler) List(c *fiber.Ctx) error {
	page := max(c.QueryInt("page", 1), 1)
	limit := min(max(c.QueryInt("limit", 10), 1), 100)

	recipes, total, err := h.repo.List(c.Context(), c.Query("q"), limit, (page-1)*limit)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"meta": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
		"data": recipes,
	})

}

func (h *RecipeHandler) GetByID(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}

	recipe, err := h.repo.GetByID(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"data": recipe})
}

func (h *RecipeHandler) UpdateRecipe(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}

	var in model.RecipeUpdateInput
	if err := c.BodyParser(&in); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "body request tidak valid")
	}
	if msg := in.ValidateUpdate(); msg != "" {
		return fiber.NewError(fiber.StatusUnprocessableEntity, msg)
	}

	recipe, err := h.repo.Update(c.Context(), id, in)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "resep tidak ditemukan"})
		}
		return err
	}
	return c.JSON(fiber.Map{"data": recipe})
}

func (h *RecipeHandler) Delete(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}

	errs := h.repo.Delete(c.Context(), id, auth.UserID(c))
	if errs != nil {
		return errs
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	msg := "terjadi kesalahan pada server"

	var fe *fiber.Error
	switch {
	case errors.Is(err, repository.ErrNotFound):
		code, msg = fiber.StatusNotFound, err.Error()
	case errors.Is(err, repository.ErrEmailTaken):
		code, msg = fiber.StatusConflict, err.Error()
	case errors.As(err, &fe):
		code, msg = fe.Code, fe.Message
	default:
		log.Printf("internal error: %v", err)
	}
	return c.Status(code).JSON(fiber.Map{"error": msg})
}

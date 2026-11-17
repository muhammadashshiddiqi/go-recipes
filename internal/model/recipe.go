package model

import (
	"strings"
	"time"
)

type Recipe struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Ingredients     []string  `json:"ingredients"`
	Steps           []string  `json:"steps"`
	CookTimeMinutes int32     `json:"cook_time_minutes"`
	Servings        int32     `json:"servings"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type RecipeInput struct {
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Ingredients     []string `json:"ingredients"`
	Steps           []string `json:"steps"`
	CookTimeMinutes int32    `json:"cook_time_minutes"`
	Servings        int32    `json:"servings"`
}

func (in *RecipeInput) Validate() string {
	in.Title = strings.TrimSpace(in.Title)
	in.Description = strings.TrimSpace(in.Description)
	switch {
	case in.Title == "":
		return "title is required"
	case len(in.Title) > 150:
		return "title maximum 150 character"
	case len(in.Ingredients) == 0:
		return "ingredients minimal 1 item"
	case len(in.Steps) == 0:
		return "steps minimal 1 steps"
	case in.CookTimeMinutes < 0:
		return "cook_time_minutes not negative"
	}

	if in.Servings <= 0 {
		in.Servings = 1
	}
	return ""
}

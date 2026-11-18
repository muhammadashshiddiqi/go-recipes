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
type RecipeUpdateInput struct {
	Title           *string  `json:"title"`
	Description     *string  `json:"description"`
	Ingredients     []string `json:"ingredients"`
	Steps           []string `json:"steps"`
	CookTimeMinutes *int32   `json:"cook_time_minutes"`
	Servings        *int32   `json:"servings"`
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

func (in *RecipeUpdateInput) ValidateUpdate() string {
	if in.IsEmpty() {
		return "minimal satu field harus diisi"
	}
	if in.Title != nil {
		t := strings.TrimSpace(*in.Title)
		in.Title = &t
		if t == "" {
			return "title tidak boleh kosong"
		}
		if len(t) > 150 {
			return "title maksimal 150 karakter"
		}
	}
	if in.Description != nil {
		d := strings.TrimSpace(*in.Description)
		in.Description = &d
	}
	if in.Ingredients != nil && len(in.Ingredients) == 0 {
		return "ingredients minimal 1 item"
	}
	if in.Steps != nil && len(in.Steps) == 0 {
		return "steps minimal 1 langkah"
	}
	if in.CookTimeMinutes != nil && *in.CookTimeMinutes < 0 {
		return "cook_time_minutes tidak boleh negatif"
	}
	if in.Servings != nil && *in.Servings <= 0 {
		return "servings harus lebih dari 0"
	}
	return ""
}

func (in *RecipeUpdateInput) IsEmpty() bool {
	return in.Title == nil && in.Description == nil && in.Ingredients == nil &&
		in.Steps == nil && in.CookTimeMinutes == nil && in.Servings == nil
}

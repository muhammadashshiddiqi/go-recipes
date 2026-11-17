package model

import (
	"strings"
	"time"
)

type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (in *RegisterInput) Validate() string {
	in.Name = strings.TrimSpace(in.Name)
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))
	switch {
	case in.Name == "":
		return "name is required"
	case !strings.Contains(in.Email, "@"):
		return "email not valid"
	case len(in.Password) < 8:
		return "password minimum 8 character"
	case len(in.Password) > 72:
		return "password maximum 72 character"
	}
	return ""
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

package handler

import (
	"errors"
	"go-recipe/internal/auth"
	"go-recipe/internal/model"
	"go-recipe/internal/repository"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	users  *repository.UserRepository
	secret string
	ttl    time.Duration
}

func NewAuthHandler(users *repository.UserRepository, secret string, ttl time.Duration) *AuthHandler {
	return &AuthHandler{users: users, secret: secret, ttl: ttl}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var in model.RegisterInput
	if err := c.BodyParser(&in); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "body request tidak valid")
	}

	if msg := in.Validate(); msg != "" {
		return fiber.NewError(fiber.StatusUnprocessableEntity, msg)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user, err := h.users.Create(c.Context(), in.Name, in.Email, string(hash))

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": user})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var in model.LoginInput

	err := c.BodyParser(&in)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "body request tidak valid")
	}
	email := strings.ToLower(strings.TrimSpace(in.Email))

	invalid := fiber.NewError(fiber.StatusUnauthorized, "email atau password salah")

	user, err := h.users.GetByEmail(c.Context(), email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return invalid
		}
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)) != nil {
		return invalid
	}

	token, err := auth.GenerateToken(h.secret, user.ID, h.ttl)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"token":      token,
			"token_type": "Bearer",
			"user":       user,
		},
	})
}

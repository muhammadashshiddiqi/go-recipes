package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

const userIDKey = "user_id"

func Protected(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		parts := strings.SplitN(c.Get("Authorization"), " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return fiber.NewError(fiber.StatusUnauthorized, "token not found")
		}
		claims, err := ParseToken(secret, parts[1])

		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "token tidak valid atau kadaluwarsa")
		}

		c.Locals(userIDKey, claims.UserID)
		return c.Next()
	}
}

func UserID(c *fiber.Ctx) int64 {
	id, _ := c.Locals(userIDKey).(int64)
	return id
}

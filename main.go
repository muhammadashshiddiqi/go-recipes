package main

import (
	"context"
	"go-recipe/internal/auth"
	"go-recipe/internal/config"
	"go-recipe/internal/database"
	"go-recipe/internal/handler"
	"go-recipe/internal/repository"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	cfg := config.Load()
	if len(cfg.JWTSecret) < 16 {
		log.Fatal("JWT_SECRET wajib diisi (min 16 karakter)")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	cancel()

	if err != nil {
		log.Fatalf("failed connect database : %v", err)
	}
	defer pool.Close()

	userRepo := repository.NewUserRepository(pool)
	recipeRepo := repository.NewRecipeRepository(pool)
	authHandler := handler.NewAuthHandler(userRepo, cfg.JWTSecret, cfg.JWTTL)
	recipeHandler := handler.NewRecipeHandler(recipeRepo)

	app := fiber.New(fiber.Config{
		AppName:      "Recipe API",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorHandler: handler.ErrorHandler,
	})
	app.Use(recover.New(), logger.New(), cors.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		if err := pool.Ping(c.Context()); err != nil {
			return fiber.NewError(fiber.StatusServiceUnavailable, "database tidak tersedia")
		}

		return c.JSON(fiber.Map{"status": "ok"})
	})

	api := app.Group("/api/v1")
	api.Post("/auth/register", authHandler.Register)
	api.Post("/auth/login", authHandler.Login)
	api.Get("/recipes", recipeHandler.GetByID)

	protected := auth.Protected(cfg.JWTSecret)
	api.Post("/recipes", protected, recipeHandler.Create)
	api.Put("recipes/:id", protected, recipeHandler.Update)
	api.Delete("/recipes/:id", protected, recipeHandler.Delete)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		_ = app.Shutdown()
	}()

	log.Fatal(app.Listen(":" + cfg.Port))
}

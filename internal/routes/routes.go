package routes

import (
	"mca-action-convertor/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(
	app *fiber.App,
	handler *handlers.Handler,
) {
	// API group
	api := app.Group("/api")

	// Version group
	v1 := api.Group("/v1")

	// Health check route
	v1.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "API is up and running",
		})
	})

	// Converter endpoints
	converter := v1.Group("/convert")
	converter.Post("/", handler.ConvertJSON)

	// Not found handler
	app.Use(func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusNotFound, "Endpoint not found")
	})
}

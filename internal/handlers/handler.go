package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"mca-action-convertor/internal/usecase"
)

// Handler handles HTTP requests
type Handler struct {
	converterUseCase *usecase.QueryConverterUseCase
	logger           *zap.Logger
}

func NewHandler(converterUseCase *usecase.QueryConverterUseCase, logger *zap.Logger) *Handler {
	return &Handler{
		converterUseCase: converterUseCase,
		logger:           logger.With(zap.String("handler", "handler")),
	}
}

// ConvertJSON handles POST /convert/json endpoint
func (h *Handler) ConvertJSON(c *fiber.Ctx) error {
	// Read raw JSON from request body
	body := c.Body()
	if len(body) == 0 {
		h.logger.Warn("Invalid request body")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Convert JSON to SQL
	sqlMap, err := h.converterUseCase.ConvertJSONToSQL(string(body))
	if err != nil {
		h.logger.Warn("Failed to convert JSON", zap.Error(err))
		return fiber.NewError(fiber.StatusBadRequest, "Failed to convert JSON")
	}

	// Return response
	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"queries": sqlMap,
		},
	})
}

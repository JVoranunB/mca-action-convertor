// internal/handlers/error_handler.go
package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// Logger is a package-level variable that should be set from main
var Logger *zap.Logger

// InitLogger sets the logger for the handlers package
func InitLogger(logger *zap.Logger) {
	Logger = logger.With(zap.String("component", "errorHandler"))
}

// ErrorHandler is a custom error handler for Fiber
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Default status code is 500
	code := fiber.StatusInternalServerError

	// Retrieve the custom status code if it's a fiber.*Error
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	// Handle GORM record not found error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		code = fiber.StatusNotFound
	}

	// Get request information for context
	path := c.Path()
	method := c.Method()
	ip := c.IP()

	// Log the error with context
	if Logger != nil {
		Logger.Error("Request error",
			zap.String("method", method),
			zap.String("path", path),
			zap.String("ip", ip),
			zap.Int("status", code),
			zap.Error(err))
	}

	// Return the error as JSON
	return c.Status(code).JSON(ErrorResponse{
		Error:   "Error processing request",
		Message: err.Error(),
	})
}

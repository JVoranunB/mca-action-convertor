// cmd/api/main.go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"mca-bigQuery/internal/adapter/jsonparser"
	"mca-bigQuery/internal/adapter/sqlbuilder"
	"mca-bigQuery/internal/handlers"
	"mca-bigQuery/internal/infrastructure/config"
	"mca-bigQuery/internal/infrastructure/logger"
	"mca-bigQuery/internal/repository"
	"mca-bigQuery/internal/routes"
	"mca-bigQuery/internal/usecase"

	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Initialize logger
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	log, err := logger.Initialize(env)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer log.Sync()

	sugar := log.Sugar()
	sugar.Info("Marketing Automation API Starting...")

	// Initialize error handler logger
	handlers.InitLogger(log)

	// Initialize repositories
	parser := jsonparser.NewParser()
	repo := repository.NewQueryRepository(parser)
	sqlBuilder := sqlbuilder.NewSQLBuilder()
	converter := usecase.NewQueryConverterUseCase(repo, sqlBuilder)
	handler := handlers.NewHandler(converter, log)

	// Create a new Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: handlers.ErrorHandler,
	})

	// Middleware
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: log,
	}))
	app.Use(recover.New())
	app.Use(cors.New())

	// Setup routes
	routes.SetupRoutes(app, handler)

	// Get port from environment variable or use default
	port := config.GetEnv("PORT", "8080")

	// Start server in a goroutine
	go func() {
		sugar.Infof("Server starting on port %s", port)
		if err := app.Listen(":" + port); err != nil {
			sugar.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	sugar.Info("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		sugar.Fatalf("Server shutdown failed: %v", err)
	}

	sugar.Info("Server exited properly")
}

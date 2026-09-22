package main

import (
	"log"
	"os"

	"eposyandu-backend/config"
	"eposyandu-backend/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// 1. Inisialisasi Koneksi Database
	config.ConnectDB()

	// 2. Inisialisasi App Fiber
	app := fiber.New(fiber.Config{
		AppName: "e-Posyandu API v1.0",
	})

	// 3. Global Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// 4. Setup API Routes
	routes.SetupRoutes(app)

	// 5. Health Check Endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "e-Posyandu Backend Service is running smoothly",
		})
	})

	// 6. Jalankan Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server berjalan di http://localhost:%s", port)
	log.Fatal(app.Listen(":" + port))
}
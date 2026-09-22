package main

import (
	"log"
	"os"

	"eposyandu-backend/config"
	"eposyandu-backend/database"
	"eposyandu-backend/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load File .env (Mencegah error Access Denied MySQL)
	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan environment variable bawaan/default OS")
	}

	// 2. Inisialisasi Koneksi Database
	config.ConnectDB()

	// 3. Jalankan Database Seeder (Membuat Akun Admin, Kader, Ortu bawaan)
	database.SeedUser()

	// 4. Inisialisasi App Fiber
	app := fiber.New(fiber.Config{
		AppName: "e-Posyandu API v1.0",
	})

	// 5. Global Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// 6. Setup API Routes
	routes.SetupRoutes(app)

	// 7. Health Check Endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "e-Posyandu Backend Service is running smoothly",
		})
	})

	// 8. Jalankan Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server berjalan di http://localhost:%s", port)
	log.Fatal(app.Listen(":" + port))
}
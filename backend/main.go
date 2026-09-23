package main

import (
	"log"
	"os"

	"eposyandu-backend/config"
	"eposyandu-backend/database"
	_ "eposyandu-backend/docs" // Blank import untuk menginisialisasi Swagger docs
	"eposyandu-backend/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

// @title           e-Posyandu API Documentation
// @version         1.0
// @description     API Service Backend untuk Sistem Informasi e-Posyandu berbasis Go Fiber.
// @termsOfService  http://swagger.io/terms/

// @contact.name    Tim Pengembang e-Posyandu
// @contact.email   admin@eposyandu.local

// @license.name    MIT
// @license.url     https://opensource.org/licenses/MIT

// @host            localhost:3000
// @BasePath        /api/v1

// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
// @description                Masukkan token JWT dengan format: Bearer <token>
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

	// 6. Endpoint Swagger UI
	app.Get("/swagger/*", swagger.HandlerDefault)

	// 7. Health Check Endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "e-Posyandu Backend Service is running smoothly",
		})
	})

	// 8. Setup API Routes
	routes.SetupRoutes(app)

	// 9. Jalankan Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server berjalan di http://localhost:%s", port)
	log.Printf("Dokumentasi Swagger dapat diakses di http://localhost:%s/swagger/index.html", port)
	log.Fatal(app.Listen(":" + port))
}
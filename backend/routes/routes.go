package routes

import (
	"eposyandu-backend/controllers"
	"eposyandu-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	// Public Routes
	auth := api.Group("/auth")
	auth.Post("/register", controllers.Register)
	auth.Post("/login", controllers.Login)

	// Protected Routes (Perlu JWT)
	protected := api.Group("", middleware.Protected())

	// Profile
	protected.Get("/user/profile", controllers.GetProfile)

	// Modul Anak
	protected.Post("/anak", middleware.AuthorizeRoles("admin", "nakes", "kader", "ibu"), controllers.CreateAnak)
	protected.Get("/anak", middleware.AuthorizeRoles("admin", "nakes", "kader"), controllers.GetAllAnak)
	protected.Get("/anak/my-children", middleware.AuthorizeRoles("ibu"), controllers.GetAnakByIbu)

	// Modul KMS (Penimbangan)
	protected.Post("/kms", middleware.AuthorizeRoles("admin", "nakes", "kader"), controllers.AddPenimbangan)
	protected.Get("/kms/anak/:anak_id", controllers.GetRiwayatKMSByAnak)

	// Modul Imunisasi
	protected.Post("/imunisasi", middleware.AuthorizeRoles("admin", "nakes", "kader"), controllers.AddImunisasi)
	protected.Get("/imunisasi/anak/:anak_id", controllers.GetImunisasiByAnak)

	// Modul Jadwal Posyandu
	protected.Post("/jadwal", middleware.AuthorizeRoles("admin", "nakes", "kader"), controllers.CreateJadwal)
	protected.Get("/jadwal", controllers.GetJadwal)

	// Modul Dashboard (Nakes/Kader)
	protected.Get("/dashboard/summary", middleware.AuthorizeRoles("admin", "nakes", "kader"), controllers.GetDashboardSummary)
}
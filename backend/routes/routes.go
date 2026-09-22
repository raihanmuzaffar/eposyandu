package routes

import (
	"eposyandu-backend/controllers"
	"eposyandu-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	// 1. Public Routes
	auth := api.Group("/auth")
	auth.Post("/register", controllers.Register)
	auth.Post("/login", controllers.Login)

	// 2. Protected Routes (Membutuhkan JWT Token Valid)
	protected := api.Group("", middleware.Protected())

	// Profile & User Management
	protected.Get("/user/profile", controllers.GetProfile)

	// Modul Anak
	anak := protected.Group("/anak")
	anak.Post("/", middleware.AuthorizeRoles("ADMIN", "KADER", "ORTU"), controllers.CreateAnak)
	anak.Get("/", middleware.AuthorizeRoles("ADMIN", "KADER"), controllers.GetAllAnak)
	anak.Get("/my-children", middleware.AuthorizeRoles("ORTU"), controllers.GetAnakByIbu)

	// Modul KMS (Penimbangan)
	kms := protected.Group("/kms")
	kms.Post("/", middleware.AuthorizeRoles("ADMIN", "KADER"), controllers.AddPenimbangan)
	kms.Get("/anak/:anak_id", middleware.AuthorizeRoles("ADMIN", "KADER", "ORTU"), controllers.GetRiwayatKMSByAnak)

	// Modul Imunisasi
	imunisasi := protected.Group("/imunisasi")
	imunisasi.Post("/", middleware.AuthorizeRoles("ADMIN", "KADER"), controllers.AddImunisasi)
	imunisasi.Get("/anak/:anak_id", middleware.AuthorizeRoles("ADMIN", "KADER", "ORTU"), controllers.GetImunisasiByAnak)

	// Modul Ibu Hamil & Pemeriksaan
	ibuHamil := protected.Group("/ibu-hamil")
	ibuHamil.Post("/", middleware.AuthorizeRoles("ADMIN", "KADER"), controllers.RegisterIbuHamil)
	ibuHamil.Get("/", middleware.AuthorizeRoles("ADMIN", "KADER"), controllers.GetAllIbuHamil)
	ibuHamil.Post("/pemeriksaan", middleware.AuthorizeRoles("ADMIN", "KADER"), controllers.AddPemeriksaanBumil)
	ibuHamil.Get("/pemeriksaan/:ibu_hamil_id", middleware.AuthorizeRoles("ADMIN", "KADER", "ORTU"), controllers.GetRiwayatPemeriksaanBumil)

	// Modul Jadwal Posyandu
	jadwal := protected.Group("/jadwal")
	jadwal.Post("/", middleware.AuthorizeRoles("ADMIN", "KADER"), controllers.CreateJadwal)
	jadwal.Get("/", controllers.GetJadwal) // Dapat diakses semua role terautentikasi

	// Modul Edukasi & Informasi
	edukasi := protected.Group("/edukasi")
	edukasi.Post("/", middleware.AuthorizeRoles("ADMIN", "KADER"), controllers.CreateEdukasi)
	edukasi.Get("/", controllers.GetAllEdukasi)    // Dapat diakses semua role terautentikasi
	edukasi.Get("/:id", controllers.GetEdukasiByID) // Dapat diakses semua role terautentikasi

	// Modul Dashboard
	dashboard := protected.Group("/dashboard")
	dashboard.Get("/summary", middleware.AuthorizeRoles("ADMIN", "KADER"), controllers.GetDashboardSummary)
}
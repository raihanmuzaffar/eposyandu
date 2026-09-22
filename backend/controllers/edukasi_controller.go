package controllers

import (
	"eposyandu-backend/config"
	"eposyandu-backend/models"

	"github.com/gofiber/fiber/v2"
)

type CreateEdukasiInput struct {
	Judul     string `json:"judul" validate:"required"`
	Kategori  string `json:"kategori" validate:"required"`
	Konten    string `json:"konten" validate:"required"`
	GambarURL string `json:"gambar_url"`
}

// 1. Tambah Artikel / Pengumuman Baru (ADMIN & KADER)
func CreateEdukasi(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var input CreateEdukasiInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format input tidak valid",
		})
	}

	artikel := models.Edukasi{
		Judul:     input.Judul,
		Kategori:  input.Kategori,
		Konten:    input.Konten,
		GambarURL: input.GambarURL,
		PenulisID: userID,
	}

	if err := config.DB.Create(&artikel).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal menyimpan konten edukasi",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Konten edukasi berhasil dipublikasikan",
		"data":    artikel,
	})
}

// 2. Mengambil Semua Artikel Edukasi (Publik / Semua Role)
func GetAllEdukasi(c *fiber.Ctx) error {
	var daftarEdukasi []models.Edukasi
	if err := config.DB.Order("created_at desc").Find(&daftarEdukasi).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mengambil daftar edukasi",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   daftarEdukasi,
	})
}

// 3. Mengambil Detail Artikel Berdasarkan ID
func GetEdukasiByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var artikel models.Edukasi
	if err := config.DB.First(&artikel, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Konten edukasi tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   artikel,
	})
}
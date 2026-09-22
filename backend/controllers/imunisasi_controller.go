package controllers

import (
	"time"

	"eposyandu-backend/config"
	"eposyandu-backend/models"

	"github.com/gofiber/fiber/v2"
)

type CreateImunisasiInput struct {
	AnakID        uint   `json:"anak_id" validate:"required"`
	JenisVaksin   string `json:"jenis_vaksin" validate:"required"`
	TanggalSuntik string `json:"tanggal_suntik" validate:"required"` // Format: YYYY-MM-DD
	VitaminA      string `json:"vitamin_a"`
	ObatCacing    bool   `json:"obat_cacing"`
	Catatan       string `json:"catatan"`
}

func AddImunisasi(c *fiber.Ctx) error {
	kaderID := c.Locals("user_id").(uint)

	var input CreateImunisasiInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format input tidak valid",
		})
	}

	tglSuntik, err := time.Parse("2006-01-02", input.TanggalSuntik)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format tanggal_suntik harus YYYY-MM-DD",
		})
	}

	// Pastikan data anak ada
	var anak models.Anak
	if err := config.DB.First(&anak, input.AnakID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Data anak tidak ditemukan",
		})
	}

	imunisasi := models.Imunisasi{
		AnakID:        input.AnakID,
		KaderID:       kaderID,
		JenisVaksin:   input.JenisVaksin,
		TanggalSuntik: tglSuntik,
		VitaminA:      input.VitaminA,
		ObatCacing:    input.ObatCacing,
		Catatan:       input.Catatan,
	}

	if err := config.DB.Create(&imunisasi).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal menyimpan catatan imunisasi",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Catatan imunisasi berhasil ditambahkan",
		"data":    imunisasi,
	})
}

func GetImunisasiByAnak(c *fiber.Ctx) error {
	anakID := c.Params("anak_id")

	var daftarImunisasi []models.Imunisasi
	if err := config.DB.Where("anak_id = ?", anakID).Order("tanggal_suntik desc").Find(&daftarImunisasi).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mengambil riwayat imunisasi",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   daftarImunisasi,
	})
}
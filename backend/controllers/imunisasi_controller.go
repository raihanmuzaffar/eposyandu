package controllers

import (
	"time"

	"eposyandu-backend/config"
	"eposyandu-backend/models"

	"github.com/gofiber/fiber/v2"
)

type CreateImunisasiInput struct {
	AnakID           uint   `json:"anak_id"`
	JenisImunisasi   string `json:"jenis_imunisasi"`
	TanggalPemberian string `json:"tanggal_pemberian"` // YYYY-MM-DD
	Catatan          string `json:"catatan"`
}

func AddImunisasi(c *fiber.Ctx) error {
	kaderID := c.Locals("user_id").(uint)

	var input CreateImunisasiInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format input tidak valid"})
	}

	tgl, err := time.Parse("2006-01-02", input.TanggalPemberian)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format tanggal harus YYYY-MM-DD"})
	}

	imunisasi := models.Imunisasi{
		AnakID:           input.AnakID,
		KaderID:          kaderID,
		JenisImunisasi:   input.JenisImunisasi,
		TanggalPemberian: tgl,
		Catatan:          input.Catatan,
	}

	if err := config.DB.Create(&imunisasi).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal mencatat imunisasi"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Catatan imunisasi berhasil ditambahkan",
		"data":    imunisasi,
	})
}

func GetImunisasiByAnak(c *fiber.Ctx) error {
	anakID := c.Params("anak_id")

	var list []models.Imunisasi
	if err := config.DB.Where("anak_id = ?", anakID).Order("tanggal_pemberian desc").Find(&list).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal mengambil riwayat imunisasi"})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   list,
	})
}
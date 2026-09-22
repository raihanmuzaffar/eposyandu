package controllers

import (
	"time"

	"eposyandu-backend/config"
	"eposyandu-backend/models"

	"github.com/gofiber/fiber/v2"
)

type CreateJadwalInput struct {
	PosyanduID   *uint  `json:"posyandu_id"`
	Judul        string `json:"judul" validate:"required"`
	Kegiatan     string `json:"kegiatan" validate:"required"`
	Lokasi       string `json:"lokasi" validate:"required"`
	Tanggal      string `json:"tanggal" validate:"required"` // YYYY-MM-DD
	WaktuMulai   string `json:"waktu_mulai" validate:"required"`
	WaktuSelesai string `json:"waktu_selesai" validate:"required"`
}

func CreateJadwal(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var input CreateJadwalInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format input tidak valid",
		})
	}

	tgl, err := time.Parse("2006-01-02", input.Tanggal)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format tanggal harus YYYY-MM-DD",
		})
	}

	jadwal := models.JadwalPosyandu{
		PosyanduID:   input.PosyanduID,
		Judul:        input.Judul,
		Kegiatan:     input.Kegiatan,
		Lokasi:       input.Lokasi,
		Tanggal:      tgl,
		WaktuMulai:   input.WaktuMulai,
		WaktuSelesai: input.WaktuSelesai,
		CreatedByID:  userID,
	}

	if err := config.DB.Create(&jadwal).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal membuat jadwal posyandu",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Jadwal kegiatan berhasil ditambahkan",
		"data":    jadwal,
	})
}

func GetJadwal(c *fiber.Ctx) error {
	var daftarJadwal []models.JadwalPosyandu
	if err := config.DB.Order("tanggal desc").Find(&daftarJadwal).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mengambil daftar jadwal",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   daftarJadwal,
	})
}
package controllers

import (
	"time"

	"eposyandu-backend/config"
	"eposyandu-backend/models"

	"github.com/gofiber/fiber/v2"
)

type CreateJadwalInput struct {
	PosyanduID uint   `json:"posyandu_id"`
	Judul      string `json:"judul"`
	Kegiatan   string `json:"kegiatan"`
	Tanggal    string `json:"tanggal"` // Format: YYYY-MM-DD HH:mm
	Lokasi     string `json:"lokasi"`
}

func CreateJadwal(c *fiber.Ctx) error {
	var input CreateJadwalInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format input tidak valid"})
	}

	tgl, err := time.Parse("2006-01-02 15:04", input.Tanggal)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format tanggal harus YYYY-MM-DD HH:mm"})
	}

	jadwal := models.JadwalPosyandu{
		PosyanduID: input.PosyanduID,
		Judul:      input.Judul,
		Kegiatan:   input.Kegiatan,
		Tanggal:    tgl,
		Lokasi:     input.Lokasi,
	}

	if err := config.DB.Create(&jadwal).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal membuat jadwal"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Jadwal posyandu berhasil dibuat",
		"data":    jadwal,
	})
}

func GetJadwal(c *fiber.Ctx) error {
	var jadwalList []models.JadwalPosyandu
	if err := config.DB.Where("tanggal >= ?", time.Now()).Order("tanggal asc").Find(&jadwalList).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal mengambil jadwal posyandu"})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   jadwalList,
	})
}
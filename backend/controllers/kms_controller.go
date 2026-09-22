package controllers

import (
	"time"

	"eposyandu-backend/config"
	"eposyandu-backend/models"
	"eposyandu-backend/utils"

	"github.com/gofiber/fiber/v2"
)

type CreateKMSInput struct {
	AnakID          uint    `json:"anak_id"`
	TanggalTimbang  string  `json:"tanggal_timbang"` // YYYY-MM-DD
	UsiaBulan       int     `json:"usia_bulan"`
	BeratBadanKg    float64 `json:"berat_badan_kg"`
	TinggiBadanCm   float64 `json:"tinggi_badan_cm"`
	LingkarKepalaCm float64 `json:"lingkar_kepala_cm"`
	Catatan         string  `json:"catatan"`
}

func AddPenimbangan(c *fiber.Ctx) error {
	kaderID := c.Locals("user_id").(uint)

	var input CreateKMSInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format input tidak valid"})
	}

	tglTimbang, err := time.Parse("2006-01-02", input.TanggalTimbang)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format tanggal_timbang harus YYYY-MM-DD"})
	}

	var anak models.Anak
	if err := config.DB.First(&anak, input.AnakID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Data anak tidak ditemukan"})
	}

	// Cast anak.JenisKelamin (models.JenisKelamin) ke string
	statusGizi := utils.HitungZScoreStatusGizi(input.BeratBadanKg, input.UsiaBulan, string(anak.JenisKelamin))
	statusStunting := utils.HitungZScoreStatusStunting(input.TinggiBadanCm, input.UsiaBulan, string(anak.JenisKelamin))

	kms := models.PenimbanganKMS{
		AnakID:          input.AnakID,
		KaderID:         kaderID,
		TanggalTimbang:  tglTimbang,
		UsiaBulan:       input.UsiaBulan,
		BeratBadanKg:    input.BeratBadanKg,
		TinggiBadanCm:   input.TinggiBadanCm,
		LingkarKepalaCm: input.LingkarKepalaCm,
		StatusGizi:      statusGizi,
		StatusStunting:  statusStunting,
		Catatan:         input.Catatan,
	}

	if err := config.DB.Create(&kms).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal menyimpan data penimbangan"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Catatan KMS berhasil ditambahkan dan Z-Score WHO dihitung",
		"data":    kms,
	})
}

// GetRiwayatKMSByAnak melihat histori penimbangan KMS balita
func GetRiwayatKMSByAnak(c *fiber.Ctx) error {
	anakID := c.Params("anak_id")

	var riwayat []models.PenimbanganKMS
	if err := config.DB.Where("anak_id = ?", anakID).Order("usia_bulan asc").Find(&riwayat).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal mengambil data riwayat KMS"})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   riwayat,
	})
}
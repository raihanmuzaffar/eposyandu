package controllers

import (
	"time"

	"eposyandu-backend/config"
	"eposyandu-backend/models"
	"eposyandu-backend/utils"

	"github.com/gofiber/fiber/v2"
)

type CreateKMSInput struct {
	AnakID          uint    `json:"anak_id" example:"1"`
	TanggalTimbang  string  `json:"tanggal_timbang" example:"2026-09-22"` // YYYY-MM-DD
	UsiaBulan       int     `json:"usia_bulan" example:"12"`
	BeratBadanKg    float64 `json:"berat_badan_kg" example:"9.5"`
	TinggiBadanCm   float64 `json:"tinggi_badan_cm" example:"75.0"`
	LingkarKepalaCm float64 `json:"lingkar_kepala_cm" example:"45.5"`
	Catatan         string  `json:"catatan" example:"Tumbuh kembang baik, berikan ASI dan MPASI bergizi"`
}

// AddPenimbangan godoc
// @Summary      Catat Penimbangan & Tumbuh Kembang Anak (KMS)
// @Description  Mencatat hasil penimbangan balita dan mengkalkulasi otomatis status Z-Score WHO (Status Gizi & Stunting)
// @Tags         KMS & Penimbangan
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateKMSInput true "Data Penimbangan Anak"
// @Success      201 {object} map[string]interface{} "Catatan KMS berhasil ditambahkan"
// @Failure      400 {object} map[string]interface{} "Format request atau tanggal tidak valid"
// @Failure      401 {object} map[string]interface{} "Tidak terautentikasi"
// @Failure      404 {object} map[string]interface{} "Data anak tidak ditemukan"
// @Failure      500 {object} map[string]interface{} "Gagal menyimpan data penimbangan"
// @Router       /kms [post]
func AddPenimbangan(c *fiber.Ctx) error {
	// Pengecekan ID Kader secara aman dari c.Locals untuk mencegah panic
	var kaderID uint
	if id, ok := c.Locals("user_id").(uint); ok {
		kaderID = id
	} else if id, ok := c.Locals("userID").(uint); ok {
		kaderID = id
	} else if idFloat, ok := c.Locals("user_id").(float64); ok {
		kaderID = uint(idFloat)
	} else if idFloat, ok := c.Locals("userID").(float64); ok {
		kaderID = uint(idFloat)
	}

	if kaderID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Tidak terautentikasi / ID Pengguna tidak valid",
		})
	}

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

// GetRiwayatKMSByAnak godoc
// @Summary      Riwayat KMS Balita
// @Description  Melihat histori penimbangan dan catatan tumbuh kembang balita berdasarkan ID Anak
// @Tags         KMS & Penimbangan
// @Produce      json
// @Security     BearerAuth
// @Param        anak_id path int true "ID Anak"
// @Success      200 {object} map[string]interface{} "Data riwayat KMS balita"
// @Failure      401 {object} map[string]interface{} "Tidak terautentikasi"
// @Failure      500 {object} map[string]interface{} "Gagal mengambil data riwayat KMS"
// @Router       /kms/anak/{anak_id} [get]
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
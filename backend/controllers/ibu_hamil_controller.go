package controllers

import (
	"time"

	"eposyandu-backend/config"
	"eposyandu-backend/models"

	"github.com/gofiber/fiber/v2"
)

type CreateIbuHamilInput struct {
	UserID        uint   `json:"user_id" validate:"required"`
	HPHT          string `json:"hpht" validate:"required"` // Format: YYYY-MM-DD
	GolonganDarah string `json:"golongan_darah"`
	Catatan       string `json:"catatan"`
}

type AddPemeriksaanBumilInput struct {
	IbuHamilID          uint    `json:"ibu_hamil_id" validate:"required"`
	TanggalPeriksa      string  `json:"tanggal_periksa" validate:"required"` // Format: YYYY-MM-DD
	BeratBadan          float64 `json:"berat_badan" validate:"required"`
	TekananDarahSistol  int     `json:"tekanan_darah_sistol" validate:"required"`
	TekananDarahDiastol int     `json:"tekanan_darah_diastol" validate:"required"`
	LILA                float64 `json:"lila" validate:"required"`
	JumlahTTD           int     `json:"jumlah_ttd"`
	Keluhan             string  `json:"keluhan"`
}

// 1. Pendaftaran Kehamilan Baru
func RegisterIbuHamil(c *fiber.Ctx) error {
	var input CreateIbuHamilInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format input tidak valid",
		})
	}

	tglHPHT, err := time.Parse("2006-01-02", input.HPHT)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format HPHT harus YYYY-MM-DD",
		})
	}

	// Kalkulasi Rumus Naegele untuk HPL (HPHT + 7 hari - 3 bulan + 1 tahun)
	hpl := tglHPHT.AddDate(1, -3, 7)

	ibuHamil := models.IbuHamil{
		UserID:             input.UserID,
		HPHT:               tglHPHT,
		TafsiranPersalinan: hpl,
		GolonganDarah:      input.GolonganDarah,
		Catatan:            input.Catatan,
	}

	if err := config.DB.Create(&ibuHamil).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mendaftarkan data ibu hamil",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Data kehamilan berhasil ditambahkan",
		"data":    ibuHamil,
	})
}

// 2. Mengambil Semua Data Ibu Hamil
func GetAllIbuHamil(c *fiber.Ctx) error {
	var daftarBumil []models.IbuHamil
	if err := config.DB.Preload("User").Find(&daftarBumil).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mengambil daftar ibu hamil",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   daftarBumil,
	})
}

// 3. Pencatatan Pemeriksaan Rutin Ibu Hamil
func AddPemeriksaanBumil(c *fiber.Ctx) error {
	kaderID := c.Locals("user_id").(uint)

	var input AddPemeriksaanBumilInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format input tidak valid",
		})
	}

	tglPeriksa, err := time.Parse("2006-01-02", input.TanggalPeriksa)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format tanggal_periksa harus YYYY-MM-DD",
		})
	}

	// Cek Keberadaan Data Ibu Hamil
	var bumil models.IbuHamil
	if err := config.DB.First(&bumil, input.IbuHamilID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Data ibu hamil tidak ditemukan",
		})
	}

	// Deteksi KEK (Kurang Energi Kronis) jika LILA < 23.5 cm
	if input.LILA < 23.5 {
		config.DB.Model(&bumil).Update("status_kekeringan", true)
	}

	pemeriksaan := models.PemeriksaanBumil{
		IbuHamilID:          input.IbuHamilID,
		KaderID:             kaderID,
		TanggalPeriksa:      tglPeriksa,
		BeratBadan:          input.BeratBadan,
		TekananDarahSistol:  input.TekananDarahSistol,
		TekananDarahDiastol: input.TekananDarahDiastol,
		LILA:                input.LILA,
		JumlahTTD:           input.JumlahTTD,
		Keluhan:             input.Keluhan,
	}

	if err := config.DB.Create(&pemeriksaan).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal menyimpan data pemeriksaan",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Data pemeriksaan kehamilan berhasil disimpan",
		"data":    pemeriksaan,
	})
}

// 4. Mengambil Riwayat Pemeriksaan Berdasarkan Ibu Hamil ID
func GetRiwayatPemeriksaanBumil(c *fiber.Ctx) error {
	ibuHamilID := c.Params("ibu_hamil_id")

	var riwayat []models.PemeriksaanBumil
	if err := config.DB.Where("ibu_hamil_id = ?", ibuHamilID).Order("tanggal_periksa desc").Find(&riwayat).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mengambil riwayat pemeriksaan",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   riwayat,
	})
}
package controllers

import (
	"time"

	"eposyandu-backend/config"
	"eposyandu-backend/models"

	"github.com/gofiber/fiber/v2"
)

type CreateAnakInput struct {
	IbuID          *uint   `json:"ibu_id"`
	PosyanduID     *uint   `json:"posyandu_id"`
	NamaAnak       string  `json:"nama_anak"`
	NIK            string  `json:"nik"`
	NIKIbu         string  `json:"nik_ibu"`
	TanggalLahir   string  `json:"tanggal_lahir"` // Format: YYYY-MM-DD
	JenisKelamin   string  `json:"jenis_kelamin"` // 'L' atau 'P'
	BeratLahirKg   float64 `json:"berat_lahir_kg"`
	PanjangLahirCm float64 `json:"panjang_lahir_cm"`
}

// CreateAnak (Dapat diinput oleh Kader/Nakes atau Ibu)
func CreateAnak(c *fiber.Ctx) error {
	var input CreateAnakInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format input tidak valid"})
	}

	tglLahir, err := time.Parse("2006-01-02", input.TanggalLahir)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format tanggal_lahir harus YYYY-MM-DD"})
	}

	anak := models.Anak{
		IbuID:          input.IbuID,
		PosyanduID:     input.PosyanduID,
		Nama:           input.NamaAnak,
		NamaAnak:       input.NamaAnak,
		NIK:            input.NIK,
		NIKIbu:         input.NIKIbu,
		TanggalLahir:   tglLahir,
		JenisKelamin:   models.JenisKelamin(input.JenisKelamin), // Cast ke custom type models.JenisKelamin
		BeratLahirKg:   input.BeratLahirKg,
		PanjangLahirCm: input.PanjangLahirCm,
	}

	if err := config.DB.Create(&anak).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal menyimpan data anak"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Data anak berhasil ditambahkan",
		"data":    anak,
	})
}

// GetAllAnak (Kader/Nakes melihat daftar balita di posyandu)
func GetAllAnak(c *fiber.Ctx) error {
	var anakList []models.Anak
	if err := config.DB.Preload("User").Preload("PenimbanganKMS").Find(&anakList).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal mengambil data anak"})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   anakList,
	})
}

// GetAnakByIbu (Ibu melihat daftar anak-anaknya)
func GetAnakByIbu(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var anakList []models.Anak
	if err := config.DB.Where("ibu_id = ? OR user_id = ?", userID, userID).Preload("PenimbanganKMS").Find(&anakList).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal mengambil data anak"})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   anakList,
	})
}
package controllers

import (
	"eposyandu-backend/config"
	"eposyandu-backend/models"
	"eposyandu-backend/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Nama        string `json:"nama" xml:"nama" form:"nama"`
	NamaLengkap string `json:"nama_lengkap" xml:"nama_lengkap" form:"nama_lengkap"`
	NIK         string `json:"nik" xml:"nik" form:"nik"`
	Email       string `json:"email" xml:"email" form:"email"`
	Password    string `json:"password" xml:"password" form:"password"`
	NoHP        string `json:"no_hp" xml:"no_hp" form:"no_hp"`
	Role        string `json:"role" xml:"role" form:"role"`
}

type LoginInput struct {
	Email    string `json:"email" xml:"email" form:"email"`
	Password string `json:"password" xml:"password" form:"password"`
}

// Register handler untuk pendaftaran user baru
func Register(c *fiber.Ctx) error {
	var input RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Format request tidak valid",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal memproses password",
		})
	}

	// Casting string input ke custom type models.Role
	role := models.RoleOrtu
	if input.Role != "" {
		role = models.Role(input.Role)
	}

	user := models.User{
		Nama:        input.Nama,
		NamaLengkap: input.NamaLengkap,
		NIK:         input.NIK,
		Email:       input.Email,
		Password:    string(hashedPassword),
		NoHP:        input.NoHP,
		Role:        role,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Gagal meregistrasi pengguna. Email atau NIK mungkin sudah digunakan",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Registrasi berhasil",
		"user": fiber.Map{
			"id":    user.ID,
			"nama":  user.Nama,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// Login handler untuk otentikasi user
func Login(c *fiber.Ctx) error {
	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Format request tidak valid",
		})
	}

	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Email atau password salah",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Email atau password salah",
		})
	}

	token, err := utils.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal membuat token autentikasi",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Login berhasil",
		"token":   token,
		"user": fiber.Map{
			"id":           user.ID,
			"nama":         user.Nama,
			"nama_lengkap": user.NamaLengkap,
			"email":        user.Email,
			"role":         user.Role,
		},
	})
}

// GetProfile mengambil profil user aktif
func GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Tidak terautentikasi",
		})
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Pengguna tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"user": fiber.Map{
			"id":           user.ID,
			"nama":         user.Nama,
			"nama_lengkap": user.NamaLengkap,
			"email":        user.Email,
			"role":         user.Role,
			"nik":          user.NIK,
			"no_hp":        user.NoHP,
		},
	})
}
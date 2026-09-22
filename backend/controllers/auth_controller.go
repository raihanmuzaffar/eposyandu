package controllers

import (
	"eposyandu-backend/config"
	"eposyandu-backend/models"
	"eposyandu-backend/utils"

	"github.com/gofiber/fiber/v2"
)

type RegisterInput struct {
	NamaLengkap string `json:"nama_lengkap" validate:"required"`
	NIK         string `json:"nik"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=6"`
	Role        string `json:"role"`
	NoHP        string `json:"no_hp"`
	PosyanduID  *uint  `json:"posyandu_id"`
}

type LoginInput struct {
	NoHP     string `json:"no_hp"`
	Password string `json:"password"`
}

func Register(c *fiber.Ctx) error {
	var input RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format request tidak valid",
		})
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal memproses password",
		})
	}

	// Penanganan default role
	roleStr := input.Role
	if roleStr == "" {
		roleStr = string(models.RoleOrtu)
	}

	user := models.User{
		NamaLengkap: input.NamaLengkap,
		NIK:         input.NIK,
		Email:       input.Email,
		Password:    string(hashedPassword),
		Role:        models.Role(roleStr), // Konversi string ke models.Role
		NoHP:        input.NoHP,
		PosyanduID:  input.PosyanduID, // Menggunakan *uint langsung sesuai struct model
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mendaftarkan user, NIK/No. HP/Email mungkin sudah terdaftar",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Pendaftaran berhasil",
		"data": fiber.Map{
			"id":           user.ID,
			"nama_lengkap": user.NamaLengkap,
			"no_hp":        user.NoHP,
			"role":         user.Role,
		},
	})
}

func Login(c *fiber.Ctx) error {
	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format request tidak valid",
		})
	}

	var user models.User
	if err := config.DB.Where("no_hp = ?", input.NoHP).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "No HP atau Password salah",
		})
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "No HP atau Password salah",
		})
	}

	// Cast user.Role (models.Role) ke string saat dipassing ke GenerateJWT
	token, err := utils.GenerateJWT(user.ID, string(user.Role))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal membuat token autentikasi",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Login berhasil",
		"token":   token,
		"data": fiber.Map{
			"id":           user.ID,
			"nama_lengkap": user.NamaLengkap,
			"role":         user.Role,
		},
	})
}

func GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var user models.User
	if err := config.DB.Preload("Posyandu").First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "User tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   user,
	})
}
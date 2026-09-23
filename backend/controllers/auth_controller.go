package controllers

import (
	"eposyandu-backend/config"
	"eposyandu-backend/models"
	"eposyandu-backend/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Nama        string `json:"nama" xml:"nama" form:"nama" example:"Ahmad"`
	NamaLengkap string `json:"nama_lengkap" xml:"nama_lengkap" form:"nama_lengkap" example:"Ahmad Dahlan"`
	NIK         string `json:"nik" xml:"nik" form:"nik" example:"1234567890123456"`
	Email       string `json:"email" xml:"email" form:"email" example:"ahmad@example.com"`
	Password    string `json:"password" xml:"password" form:"password" example:"password123"`
	NoHP        string `json:"no_hp" xml:"no_hp" form:"no_hp" example:"081234567890"`
	Role        string `json:"role" xml:"role" form:"role" example:"ortu"`
}

type LoginInput struct {
	Email    string `json:"email" xml:"email" form:"email" example:"admin@eposyandu.local"`
	Password string `json:"password" xml:"password" form:"password" example:"admin123"`
}

// Register godoc
// @Summary      Registrasi Pengguna Baru
// @Description  Pendaftaran pengguna baru dengan role tertentu (default: ortu)
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body RegisterInput true "Data Registrasi Pengguna"
// @Success      201 {object} map[string]interface{} "Registrasi berhasil"
// @Failure      400 {object} map[string]interface{} "Format request atau data tidak valid"
// @Failure      500 {object} map[string]interface{} "Gagal memproses password"
// @Router       /auth/register [post]
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

// Login godoc
// @Summary      Login Pengguna
// @Description  Otentikasi pengguna untuk mendapatkan JWT Token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body LoginInput true "Kredensial Login"
// @Success      200 {object} map[string]interface{} "Login berhasil dan token diterbitkan"
// @Failure      400 {object} map[string]interface{} "Format request tidak valid"
// @Failure      401 {object} map[string]interface{} "Email atau password salah"
// @Failure      500 {object} map[string]interface{} "Gagal membuat token autentikasi"
// @Router       /auth/login [post]
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

// GetProfile godoc
// @Summary      Ambil Profil Pengguna
// @Description  Mengambil data profil lengkap pengguna yang sedang login berdasarkan JWT Token
// @Tags         Auth
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{} "Data profil pengguna"
// @Failure      401 {object} map[string]interface{} "Tidak terautentikasi"
// @Failure      404 {object} map[string]interface{} "Pengguna tidak ditemukan"
// @Router       /user/profile [get]
func GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		userID = c.Locals("user_id")
	}

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
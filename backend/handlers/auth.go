package handlers

import (
	"time"
	"eposyandu-backend/config"
	"eposyandu-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Nama     string      `json:"nama"`
	Email    string      `json:"email"`
	Password string      `json:"password"`
	Role     models.Role `json:"role"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c *fiber.Ctx) error {
	var input RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Input tidak valid"})
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	user := models.User{
		Nama:     input.Nama,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     input.Role,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal mendaftarkan user"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "User berhasil dibuat", "data": user})
}

func Login(c *fiber.Ctx) error {
	var input LoginInput
	var user models.User

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Input tidak valid"})
	}

	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "User tidak ditemukan"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"message": "Password salah"})
	}

	// Generate JWT Token
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("SECRET_KEY_POSYANDU"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal generate token"})
	}

	return c.JSON(fiber.Map{"token": t, "user": user})
}
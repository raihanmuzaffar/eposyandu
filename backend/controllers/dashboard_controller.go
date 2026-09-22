package controllers

import (
	"eposyandu-backend/config"
	"eposyandu-backend/models"

	"github.com/gofiber/fiber/v2"
)

func GetDashboardSummary(c *fiber.Ctx) error {
	var totalAnak int64
	var totalKader int64
	var totalOrtu int64
	var totalPenimbangan int64
	var totalStunting int64

	// Hitung total anak
	config.DB.Model(&models.Anak{}).Count(&totalAnak)

	// Hitung total kader & orang tua
	config.DB.Model(&models.User{}).Where("role = ?", models.RoleKader).Count(&totalKader)
	config.DB.Model(&models.User{}).Where("role = ?", models.RoleOrtu).Count(&totalOrtu)

	// Hitung total penimbangan
	config.DB.Model(&models.PenimbanganKMS{}).Count(&totalPenimbangan)

	// Hitung anak terindikasi stunting dari KMS
	config.DB.Model(&models.PenimbanganKMS{}).Where("status_stunting LIKE ?", "%Stunting%").Count(&totalStunting)

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"total_anak":        totalAnak,
			"total_kader":       totalKader,
			"total_ortu":        totalOrtu,
			"total_penimbangan": totalPenimbangan,
			"total_stunting":    totalStunting,
		},
	})
}
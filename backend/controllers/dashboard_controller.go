package controllers

import (
	"eposyandu-backend/config"
	"eposyandu-backend/models"

	"github.com/gofiber/fiber/v2"
)

func GetDashboardSummary(c *fiber.Ctx) error {
	var totalAnak int64
	var totalStunting int64
	var totalGiziBuruk int64

	config.DB.Model(&models.Anak{}).Count(&totalAnak)

	// Hitung dari rekapan KMS terakhir
	config.DB.Model(&models.PenimbanganKMS{}).
		Where("status_stunting IN ?", []string{"pendek", "sangat_pendek"}).
		Distinct("anak_id").
		Count(&totalStunting)

	config.DB.Model(&models.PenimbanganKMS{}).
		Where("status_gizi = ?", "gizi_buruk").
		Distinct("anak_id").
		Count(&totalGiziBuruk)

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"total_anak":      totalAnak,
			"total_stunting":  totalStunting,
			"total_gizi_buruk": totalGiziBuruk,
		},
	})
}
package database

import (
	"log"
	"time"

	"eposyandu-backend/config"
	"eposyandu-backend/models"

	"golang.org/x/crypto/bcrypt"
)

func SeedData() {
	// 1. Seed Default Users
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Gagal hashing password seeder: %v", err)
	}

	users := []models.User{
		{
			Email:    "admin@eposyandu.id",
			Password: string(hashedPassword),
			Role:     "admin",
			NoHP:     "081234567890",
		},
		{
			Email:    "nakes@eposyandu.id",
			Password: string(hashedPassword),
			Role:     "nakes",
			NoHP:     "081234567891",
		},
		{
			Email:    "kader@eposyandu.id",
			Password: string(hashedPassword),
			Role:     "kader",
			NoHP:     "081234567892",
		},
		{
			Email:    "ibu@eposyandu.id",
			Password: string(hashedPassword),
			Role:     "ibu",
			NoHP:     "081234567893",
		},
	}

	for _, user := range users {
		var count int64
		config.DB.Model(&models.User{}).Where("email = ?", user.Email).Count(&count)
		if count == 0 {
			config.DB.Create(&user)
			log.Printf("[Seeder] User berhasil dibuat: %s (%s)", user.Email, user.Role)
		}
	}

	// 2. Seed Master Jadwal Posyandu Awal
	var countJadwal int64
	config.DB.Model(&models.JadwalPosyandu{}).Count(&countJadwal)
	if countJadwal == 0 {
		jadwal := models.JadwalPosyandu{
			PosyanduID: 1,
			Judul:      "Posyandu Rutin Bulanan & Pemberian Vitamin A",
			Kegiatan:   "Penimbangan berat badan, pengukuran tinggi badan, lingkar kepala, dan pembagian Vitamin A.",
			Tanggal:    time.Now().AddDate(0, 0, 7),
			Lokasi:     "Posyandu Mawar I - Gedung Serbaguna",
		}
		config.DB.Create(&jadwal)
		log.Println("[Seeder] Jadwal posyandu sampel berhasil dibuat")
	}
}
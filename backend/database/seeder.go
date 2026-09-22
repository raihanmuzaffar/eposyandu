package database

import (
	"log"
	"time"

	"eposyandu-backend/config"
	"eposyandu-backend/models"

	"golang.org/x/crypto/bcrypt"
)

func SeedUser() {
	var count int64
	config.DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		return // Jangan seed ulang jika database sudah ada data user
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	// 1. Seed Users
	users := []models.User{
		{
			Nama:        "Admin Puskesmas",
			NamaLengkap: "Admin Puskesmas Utama",
			NIK:         "1101010101010001",
			Email:       "admin@eposyandu.id",
			Password:    string(hashedPassword),
			Role:        models.RoleAdmin,
			NoHP:        "081234567890",
		},
		{
			Nama:        "Kader Mawar",
			NamaLengkap: "Kader Posyandu Mawar",
			NIK:         "1101010101010002",
			Email:       "kader@eposyandu.id",
			Password:    string(hashedPassword),
			Role:        models.RoleKader,
			NoHP:        "081234567891",
		},
		{
			Nama:        "Ibu Srikandi",
			NamaLengkap: "Ibu Srikandi Utami",
			NIK:         "1101010101010003",
			Email:       "ortu@eposyandu.id",
			Password:    string(hashedPassword),
			Role:        models.RoleOrtu,
			NoHP:        "081234567892",
		},
	}

	for i := range users {
		config.DB.Create(&users[i])
	}

	log.Println("Seeder: Akun pengguna berhasil dibuat!")

	// 2. Seed Data Anak
	tglLahirAnak, _ := time.Parse("2006-01-02", "2025-01-15")
	userID := users[2].ID
	anak := models.Anak{
		NamaAnak:     "Budi Pratama",
		NIK:          "1101010101019999",
		JenisKelamin: "L",
		TanggalLahir: tglLahirAnak,
		NamaIbu:      "Ibu Srikandi Utami",
		Alamat:       "Jl. Mawar No. 12, RT 01/RW 02",
		UserID:       &userID,
	}
	config.DB.Create(&anak)

	// 3. Seed Riwayat KMS / Penimbangan (Disesuaikan dengan model PenimbanganKMS)
	tglTimbang1, _ := time.Parse("2006-01-02", "2026-07-10")
	tglTimbang2, _ := time.Parse("2006-01-02", "2026-08-10")
	kmsList := []models.PenimbanganKMS{
		{
			AnakID:          anak.ID,
			KaderID:         users[1].ID,
			TanggalTimbang:  tglTimbang1,
			UsiaBulan:       18,
			BeratBadanKg:    10.2,
			TinggiBadanCm:   80.5,
			LingkarKepalaCm: 45.0,
			StatusGizi:      "gizi_baik",
			StatusStunting:  "normal",
			Catatan:         "Kondisi sehat, tumbuh kembang baik",
		},
		{
			AnakID:          anak.ID,
			KaderID:         users[1].ID,
			TanggalTimbang:  tglTimbang2,
			UsiaBulan:       19,
			BeratBadanKg:    10.8,
			TinggiBadanCm:   82.0,
			LingkarKepalaCm: 45.5,
			StatusGizi:      "gizi_baik",
			StatusStunting:  "normal",
			Catatan:         "Nafsu makan baik, berikan ASI rutin",
		},
	}
	for i := range kmsList {
		config.DB.Create(&kmsList[i])
	}

	// 4. Seed Ibu Hamil
	hpht, _ := time.Parse("2006-01-02", "2026-02-10")
	hpl := hpht.AddDate(1, -3, 7) // Rumus Naegele
	ibuHamil := models.IbuHamil{
		UserID:             users[2].ID,
		HPHT:               hpht,
		TafsiranPersalinan: hpl,
		GolonganDarah:      "O",
		Catatan:            "Kehamilan anak kedua",
	}
	config.DB.Create(&ibuHamil)

	// 5. Seed Edukasi & Informasi
	edukasiList := []models.Edukasi{
		{
			Judul:     "Pentingnya Makanan Pendamping ASI (MPASI) Bergizi Seimbang",
			Kategori:  "Gizi",
			Konten:    "Memasuki usia 6 bulan, bayi membutuhkan nutrisi tambahan selain ASI untuk mendukung tumbuh kembang optimalnya...",
			GambarURL: "https://placehold.co/600x400/png?text=Edukasi+MPASI",
			PenulisID: users[0].ID,
		},
		{
			Judul:     "Pencegahan Stunting Sejak Masa Kehamilan",
			Kategori:  "Stunting",
			Konten:    "Pencegahan stunting harus dimulai sejak 1.000 Hari Pertama Kehidupan (HPK), yaitu sejak janin dalam kandungan...",
			GambarURL: "https://placehold.co/600x400/png?text=Cegah+Stunting",
			PenulisID: users[1].ID,
		},
	}
	for i := range edukasiList {
		config.DB.Create(&edukasiList[i])
	}

	log.Println("Seeder: Seluruh data dummy (Anak, KMS, Ibu Hamil, Edukasi) berhasil di-seed!")
}
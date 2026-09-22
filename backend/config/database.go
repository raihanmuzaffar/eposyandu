package config

import (
	"fmt"
	"log"
	"os"

	"eposyandu-backend/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dbUser := getEnv("DB_USER", "root")
	dbPass := getEnv("DB_PASSWORD", "")
	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "eposyandu_db")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi ke database: ", err)
	}

	// Auto Migrate Seluruh Model
	err = db.AutoMigrate(
		&models.User{},
		&models.Anak{},
		&models.PenimbanganKMS{},
		&models.Imunisasi{},
		&models.JadwalPosyandu{},
		&models.IbuHamil{},
		&models.PemeriksaanBumil{},
		&models.Edukasi{},
	)
	if err != nil {
		log.Fatal("Gagal migrasi database: ", err)
	}

	DB = db
	fmt.Println("Database berhasil terkoneksi & ter-migrasi!")
}

// Helper untuk membaca env dengan default value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}
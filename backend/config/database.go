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
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi ke database: ", err)
	}

	// Auto Migrate Tabel
	err = db.AutoMigrate(
		&models.User{},
		&models.Anak{},
		&models.Penimbangan{},
	)
	if err != nil {
		log.Fatal("Gagal migrasi database: ", err)
	}

	DB = db
	fmt.Println("Database berhasil terkoneksi & ter-migrasi!")
}
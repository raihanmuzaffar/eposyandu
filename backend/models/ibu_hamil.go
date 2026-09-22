package models

import (
	"time"

	"gorm.io/gorm"
)

type IbuHamil struct {
	gorm.Model
	UserID            uint      `gorm:"not null" json:"user_id"` // Reference ke User (Role: ORTU)
	HPHT              time.Time `gorm:"type:date;not null" json:"hpth"` // Hari Pertama Haid Terakhir
	TafsiranPersalinan time.Time `gorm:"type:date" json:"tafsiran_persalinan"` // HPL
	UsiaKehamilan     int       `json:"usia_kehamilan_minggu"` // Dalam Minggu
	GolonganDarah     string    `gorm:"type:varchar(5)" json:"golongan_darah"`
	StatusKekeringan  bool      `gorm:"default:false" json:"is_kek"` // Indikator KEK (Kurang Energi Kronis)
	Catatan           string    `gorm:"type:text" json:"catatan"`
	User              User      `gorm:"foreignKey:UserID" json:"-"`
}

type PemeriksaanBumil struct {
	gorm.Model
	IbuHamilID        uint      `gorm:"not null" json:"ibu_hamil_id"`
	KaderID           uint      `gorm:"not null" json:"kader_id"`
	TanggalPeriksa    time.Time `gorm:"type:date;not null" json:"tanggal_periksa"`
	BeratBadan        float64   `gorm:"type:decimal(5,2);not null" json:"berat_badan"`
	TekananDarahSistol int       `json:"tekanan_darah_sistol"` // e.g. 120
	TekananDarahDiastol int      `json:"tekanan_darah_diastol"` // e.g. 80
	LILA              float64   `gorm:"type:decimal(4,2);not null" json:"lila"` // Lingkar Lengan Atas (cm)
	JumlahTTD         int       `json:"jumlah_ttd"` // Tablet Tambah Darah
	Keluhan           string    `gorm:"type:text" json:"keluhan"`
}
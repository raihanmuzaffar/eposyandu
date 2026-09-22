package models

import (
	"time"
	"gorm.io/gorm"
)

type JenisKelamin string

const (
	LakiLaki  JenisKelamin = "L"
	Perempuan JenisKelamin = "P"
)

type Anak struct {
	gorm.Model
	NIK            string       `gorm:"type:varchar(16);uniqueIndex;not null" json:"nik"`
	Nama           string       `gorm:"type:varchar(100);not null" json:"nama"`
	NamaAnak       string       `gorm:"type:varchar(100)" json:"nama_anak"`
	JenisKelamin   JenisKelamin `gorm:"type:enum('L', 'P');not null" json:"jenis_kelamin"`
	TanggalLahir   time.Time    `gorm:"type:date;not null" json:"tanggal_lahir"`
	NamaIbu        string       `gorm:"type:varchar(100)" json:"nama_ibu"`
	NIKIbu         string       `gorm:"type:varchar(16)" json:"nik_ibu"`
	NamaAyah       string       `gorm:"type:varchar(100)" json:"nama_ayah"`
	Alamat         string       `gorm:"type:text" json:"alamat"`
	BeratLahirKg   float64      `gorm:"type:decimal(5,2)" json:"berat_lahir_kg"`
	PanjangLahirCm float64      `gorm:"type:decimal(5,2)" json:"panjang_lahir_cm"`
	IbuID          *uint        `json:"ibu_id"`
	PosyanduID     *uint        `json:"posyandu_id"`
	UserID         *uint        `json:"user_id"`
	User           User         `gorm:"foreignKey:UserID" json:"-"`
}
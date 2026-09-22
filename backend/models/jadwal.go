package models

import (
	"time"

	"gorm.io/gorm"
)

type JadwalPosyandu struct {
	gorm.Model
	PosyanduID *uint     `json:"posyandu_id"`
	Judul      string    `gorm:"type:varchar(150);not null" json:"judul"`
	Kegiatan   string    `gorm:"type:text;not null" json:"kegiatan"`
	Lokasi     string    `gorm:"type:varchar(255);not null" json:"lokasi"`
	Tanggal    time.Time `gorm:"type:date;not null" json:"tanggal"`
	WaktuMulai string    `gorm:"type:varchar(10);not null" json:"waktu_mulai"` // Contoh: "08:00"
	WaktuSelesai string  `gorm:"type:varchar(10);not null" json:"waktu_selesai"` // Contoh: "12:00"
	CreatedByID uint     `gorm:"not null" json:"created_by_id"`
}
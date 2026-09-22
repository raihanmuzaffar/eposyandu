package models

import (
	"time"

	"gorm.io/gorm"
)

type JadwalPosyandu struct {
	gorm.Model
	PosyanduID uint      `gorm:"not null" json:"posyandu_id"`
	Judul      string    `gorm:"type:varchar(100);not null" json:"judul"`
	Kegiatan   string    `gorm:"type:text" json:"kegiatan"`
	Tanggal    time.Time `gorm:"not null" json:"tanggal"`
	Lokasi     string    `gorm:"type:varchar(255)" json:"lokasi"`
}
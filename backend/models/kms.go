package models

import (
	"time"
	"gorm.io/gorm"
)

type PenimbanganKMS struct {
	gorm.Model
	AnakID          uint      `gorm:"not null" json:"anak_id"`
	KaderID         uint      `gorm:"not null" json:"kader_id"`
	TanggalTimbang  time.Time `gorm:"type:date;not null" json:"tanggal_timbang"`
	UsiaBulan       int       `gorm:"not null" json:"usia_bulan"`
	BeratBadanKg    float64   `gorm:"type:decimal(4,2);not null" json:"berat_badan_kg"`
	TinggiBadanCm   float64   `gorm:"type:decimal(4,2);not null" json:"tinggi_badan_cm"`
	LingkarKepalaCm float64   `gorm:"type:decimal(4,2)" json:"lingkar_kepala_cm"`
	StatusGizi      string    `gorm:"type:enum('gizi_buruk','gizi_kurang','gizi_baik','berisiko_gizi_lebih','gizi_lebih','obesitas');default:'gizi_baik'" json:"status_gizi"`
	StatusStunting  string    `gorm:"type:enum('sangat_pendek','pendek','normal','tinggi');default:'normal'" json:"status_stunting"`
	Catatan         string    `gorm:"type:text" json:"catatan"`
}
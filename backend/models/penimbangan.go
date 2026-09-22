package models

import (
	"time"
	"gorm.io/gorm"
)

type Penimbangan struct {
	gorm.Model
	AnakID            uint      `gorm:"not null" json:"anak_id"`
	Anak              Anak      `gorm:"foreignKey:AnakID" json:"anak,omitempty"`
	TanggalTimbang    time.Time `gorm:"type:date;not null" json:"tanggal_timbang"`
	BeratBadan        float64   `gorm:"type:decimal(5,2);not null" json:"berat_badan"` // dalam kg
	TinggiBadan       float64   `gorm:"type:decimal(5,2);not null" json:"tinggi_badan"` // dalam cm
	LingkarKepala     float64   `gorm:"type:decimal(5,2)" json:"lingkar_kepala"` // dalam cm
	Catatan           string    `gorm:"type:text" json:"catatan"`
	StatusGizi        string    `gorm:"type:varchar(50)" json:"status_gizi"` // Hasil Kalkulasi Z-Score WHO
	KaderID           uint      `gorm:"not null" json:"kader_id"`
}
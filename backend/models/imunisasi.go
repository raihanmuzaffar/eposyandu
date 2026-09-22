package models

import (
	"time"

	"gorm.io/gorm"
)

type Imunisasi struct {
	gorm.Model
	AnakID         uint      `gorm:"not null" json:"anak_id"`
	KaderID        uint      `gorm:"not null" json:"kader_id"`
	JenisVaksin    string    `gorm:"type:varchar(100);not null" json:"jenis_vaksin"` // Contoh: HB-0, BCG, Polio 1, DPT 1, Campak
	TanggalSuntik  time.Time `gorm:"type:date;not null" json:"tanggal_suntik"`
	VitaminA       string    `gorm:"type:varchar(20)" json:"vitamin_a"`             // Biru, Merah, atau Kosong
	ObatCacing     bool      `gorm:"default:false" json:"obat_cacing"`
	Catatan        string    `gorm:"type:text" json:"catatan"`
	Anak           Anak      `gorm:"foreignKey:AnakID" json:"-"`
}
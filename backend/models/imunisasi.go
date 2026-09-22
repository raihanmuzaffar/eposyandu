package models

import (
	"time"

	"gorm.io/gorm"
)

type Imunisasi struct {
	gorm.Model
	AnakID          uint      `gorm:"not null" json:"anak_id"`
	KaderID         uint      `gorm:"not null" json:"kader_id"`
	JenisImunisasi  string    `gorm:"type:varchar(50);not null" json:"jenis_imunisasi"` // 'BCG', 'DPT-1', 'Polio-1', 'Vitamin A Merah', 'Obat Cacing'
	TanggalPemberian time.Time `gorm:"not null" json:"tanggal_pemberian"`
	Catatan         string    `gorm:"type:text" json:"catatan"`
}
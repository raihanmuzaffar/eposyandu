package models

import (
	"gorm.io/gorm"
)

type Edukasi struct {
	gorm.Model
	Judul     string `gorm:"type:varchar(200);not null" json:"judul"`
	Kategori  string `gorm:"type:varchar(50);not null" json:"kategori"` // Contoh: Stunting, Gizi, Imunisasi, Ibu Hamil
	Konten    string `gorm:"type:text;not null" json:"konten"`
	GambarURL string `gorm:"type:varchar(255)" json:"gambar_url"`
	PenulisID uint   `gorm:"not null" json:"penulis_id"`
	Penulis   User   `gorm:"foreignKey:PenulisID" json:"-"`
}
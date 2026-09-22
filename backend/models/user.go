package models

import "gorm.io/gorm"

type Role string

const (
	RoleAdmin Role = "ADMIN"
	RoleKader Role = "KADER"
	RoleOrtu  Role = "ORTU"
)

type User struct {
	gorm.Model
	Nama        string `gorm:"type:varchar(100);not null" json:"nama"`
	NamaLengkap string `gorm:"type:varchar(100)" json:"nama_lengkap"`
	NIK         string `gorm:"type:varchar(16);uniqueIndex" json:"nik"`
	Email       string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password    string `gorm:"type:varchar(255);not null" json:"-"`
	NoHP        string `gorm:"type:varchar(20)" json:"no_hp"`
	Role        Role   `gorm:"type:enum('ADMIN', 'KADER', 'ORTU');default:'ORTU'" json:"role"`
	PosyanduID  *uint  `json:"posyandu_id"`
}
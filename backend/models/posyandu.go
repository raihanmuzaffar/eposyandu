package models

import "gorm.io/gorm"

type Posyandu struct {
	gorm.Model
	NamaPosyandu  string `gorm:"type:varchar(100);not null" json:"nama_posyandu"`
	Alamat        string `gorm:"type:text;not null" json:"alamat"`
	RtRw          string `gorm:"type:varchar(20)" json:"rt_rw"`
	DesaKelurahan string `gorm:"type:varchar(100);not null" json:"desa_kelurahan"`
	Kecamatan     string `gorm:"type:varchar(100);not null" json:"kecamatan"`
	Users         []User `gorm:"foreignKey:PosyanduID" json:"users,omitempty"`
}
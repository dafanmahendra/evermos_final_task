package models

import "time"

type User struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Nama         string `json:"nama"`
	KataSandi    string `json:"kata_sandi,omitempty"`
	NoTelp       string `gorm:"unique;column:no_telp" json:"no_telp"`
	TanggalLahir string `json:"tanggal_Lahir"`
	Pekerjaan    string `json:"pekerjaan"`
	Email        string `gorm:"unique" json:"email"`
	IdProvinsi   string `json:"id_provinsi"`
	IdKota       string `json:"id_kota"`
	IsAdmin      bool   `gorm:"default:false" json:"isAdmin"` // Default user biasa

	// Relasi:  1 User punya 1 toko
	Toko Toko `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"toko,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

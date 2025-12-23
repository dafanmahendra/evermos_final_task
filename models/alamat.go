package models

import "time"

type Alamat struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `json:"user_id"` // Foreign key ke User!!!!!

	JudulAlamat  string `json:"judul_alamat"`
	NamaPenerima string `json:"nama_penerima"`
	NoTelp       string `json:"no_telp"`
	DetailAlamat string `json:"detail_alamat"`

	// API DARI EMSIFA
	ProvinsiID string `json:"provinsi_id"`
	KotaID     string `json:"kota_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

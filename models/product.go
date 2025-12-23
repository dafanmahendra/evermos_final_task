package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID            uint    `gorm:"primaryKey" json:"id"`
	NamaProduk    string  `json:"nama_produk"`
	Slug          string  `json:"slug"`
	HargaReseller float64 `json:"harga_reseller"` // Kita pake float biar aman buat duit
	HargaKonsumen float64 `json:"harga_konsumen"` // Kita pake float biar aman buat duit
	Stok          int     `json:"stok"`
	Deskripsi     string  `gorm:"type:text" json:"deskripsi"`

	// Relasi ke Toko (Barang ini punya Toko siapa?)
	TokoID uint `json:"toko_id"`
	Toko   Toko `gorm:"foreignKey:TokoID" json:"toko"` // Disini kita pake struct 'Toko'

	// --- Category Dicomment Dulu Biar Gak Error ---
	// CategoryID    uint      `json:"category_id"`
	// Category      Category  `gorm:"foreignKey:CategoryID" json:"category"`

	// Relasi ke Foto (One to Many)
	FotoProduk []FotoProduk `gorm:"foreignKey:ProductID" json:"foto_produk"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// Struct tambahan buat nyimpen URL foto
type FotoProduk struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	ProductID uint   `json:"product_id"`
	Url       string `json:"url"`
}

package models

import "time"

type Product struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	NamaProduk    string    `json:"nama_produk"`
	Slug          string    `json:"slug"`
	HargaReseller string    `json:"harga_reseller"`
	HargaKonsumen string    `json:"harga_konsumen"`
	Stok          int       `json:"stok"`
	Deskripsi     string    `gorm:"type:text" json:"deskripsi"`
	
	// Relasi ke Toko (Barang ini punya Toko siapa?)
	TokoID        uint      `json:"toko_id"`
	Toko          Toko      `gorm:"foreignKey:TokoID" json:"toko"`

	// Relasi ke Category (Barang ini jenisnya apa?)
	CategoryID    uint      `json:"category_id"`
	Category      Category  `gorm:"foreignKey:CategoryID" json:"category"`

	// Foto (Nanti buat fitur upload)
	FotoProduk    []FotoProduk `gorm:"foreignKey:ProductID" json:"foto_produk"`

	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Struct tambahan buat nyimpen URL foto (One Product has Many Photos)
type FotoProduk struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	ProductID uint   `json:"product_id"`
	Url       string `json:"url"`
}
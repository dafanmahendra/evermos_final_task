package models

import (
	"time"

	"gorm.io/gorm"
)

// TRX (Header Transaksi)
type Trx struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	UserID      uint    `json:"user_id"`
	AlamatID    uint    `json:"alamat_id"`
	Alamat      Alamat  `gorm:"foreignKey:AlamatID" json:"alamat"` 

	HargaTotal  float64 `json:"harga_total"`
	KodeInvoice string  `json:"kode_invoice"`
	Status      string  `json:"status"` // UNPAID, PAID, DONE

	// Relasi: Satu Transaksi punya banyak barang
	DetailTrx []DetailTrx `gorm:"foreignKey:TrxID" json:"detail_trx"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// DETAIL TRX (List barang di keranjang)
type DetailTrx struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	TrxID     uint `json:"trx_id"`

	// hubungin ke LogProduk (Snapshot), bukan Produk asli
	LogProdukID uint      `json:"log_produk_id"`
	LogProduk   LogProduk `gorm:"foreignKey:LogProdukID" json:"log_produk"`

	Kuantitas  int     `json:"kuantitas"`
	HargaTotal float64 `json:"harga_total"` //  UBAH JADI FLOAT
}

// LOG PRODUK (Snapshot Barang saat dibeli)
type LogProduk struct {
	ID           uint `gorm:"primaryKey" json:"id"`
	ProdukAsliID uint `json:"produk_asli_id"` // ID referensi ke tabel products asli

	NamaProduk    string  `json:"nama_produk"`
	FotoProduk    string  `json:"foto_produk"`
	HargaReseller float64 `json:"harga_reseller"`
	HargaKonsumen float64 `json:"harga_konsumen"`
	Deskripsi     string  `json:"deskripsi"`

	TokoID     uint `json:"toko_id"`
	Toko       Toko `gorm:"foreignKey:TokoID" json:"toko"`

	// CategoryID uint `json:"category_id"` // matiin dulu kalau belum ada tabel Category

	CreatedAt time.Time `json:"created_at"`
}
package controllers

import (
	"strconv"
	"strings"
	"time"

	"github.com/dafanmahendra/evermos-backend/database"
	"github.com/dafanmahendra/evermos-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// CreateProduct - Upload Barang Dagangan
func CreateProduct(c *fiber.Ctx) error {
	// 1. Ambil ID User dari Token (Hasil kerja Satpam)
	userId := c.Locals("user_id") // Ingat, ini interface{}

	// 2. Cari Toko milik User ini
	var toko models.Toko
	if err := database.DB.Where("user_id = ?", userId).First(&toko).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Lo belum punya toko bro! Bikin dulu."})
	}

	// 3. Parsing Input Produk
	var input models.Product
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Input error"})
	}

	// 4. Set Data Otomatis (Business Logic)
	input.TokoID = toko.ID // <-- INI KUNCINYA (Auto-Assign Toko)

	// Bikin Slug sederhana (HP Samsung -> hp-samsung-123)
	input.Slug = strings.ReplaceAll(strings.ToLower(input.NamaProduk), " ", "-") + "-" + strconv.Itoa(int(time.Now().Unix()))

	// 5. Simpan ke Database
	if err := database.DB.Create(&input).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal upload produk"})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Barang berhasil dijual!",
		"data":    input,
	})
}

// GetAllProducts - Buat Katalog Belanja (Bisa difilter)
func GetAllProducts(c *fiber.Ctx) error {
	var products []models.Product

	// Preload: Biar pas diambil, data Toko & Kategorinya ikut kebawa
	// Select: Ambil field public aja, jangan bawa password user dll
	query := database.DB.Preload("Toko").Preload("Category")

	// Fitur Search (Opsional - Bonus Nilai)
	search := c.Query("search")
	if search != "" {
		query = query.Where("nama_produk LIKE ?", "%"+search+"%")
	}

	query.Find(&products)

	return c.JSON(fiber.Map{
		"data": products,
	})
}

// GetProductDetail - Liat 1 Barang
func GetProductDetail(c *fiber.Ctx) error {
	productId := c.Params("id")
	var product models.Product

	if err := database.DB.Preload("Toko").Preload("Category").First(&product, productId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{"message": "Barang gak ketemu"})
		}
		return c.Status(500).JSON(fiber.Map{"message": "Error database"})
	}

	return c.JSON(fiber.Map{"data": product})
}

// UpdateProduct - Update Barang (Safe Version)
func UpdateProduct(c *fiber.Ctx) error {
	productId := c.Params("id")

	// AMAN: Konversi interface ke float64 dulu (bawaan JWT), baru ke uint
	userLocals := c.Locals("user_id")
	var userId uint
	if v, ok := userLocals.(float64); ok {
		userId = uint(v)
	} else if v, ok := userLocals.(uint); ok {
		userId = v
	} else {
		return c.Status(401).JSON(fiber.Map{"message": "User ID corrupt"})
	}

	// 1. Cek Toko User
	var toko models.Toko
	if err := database.DB.Where("user_id = ?", userId).First(&toko).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Toko tidak ditemukan / Belum bikin toko"})
	}

	// 2. Cek Barang
	var product models.Product
	if err := database.DB.First(&product, productId).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Barang gak ketemu"})
	}

	// 3. SECURITY CHECK: IDOR Protection
	if product.TokoID != toko.ID {
		return c.Status(403).JSON(fiber.Map{"message": "Dih, mau ngedit barang orang ya? Gak boleh!"})
	}

	// 4. Update Data
	var input models.Product
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Input error"})
	}

	// Update Partial
	database.DB.Model(&product).Updates(models.Product{
		NamaProduk:    input.NamaProduk,
		HargaReseller: input.HargaReseller,
		HargaKonsumen: input.HargaKonsumen,
		Stok:          input.Stok,
		Deskripsi:     input.Deskripsi,
		CategoryID:    input.CategoryID,
	})

	return c.JSON(fiber.Map{"message": "Barang berhasil diupdate", "data": product})
}

// DeleteProduct - Hapus Barang (Safe Version)
func DeleteProduct(c *fiber.Ctx) error {
	productId := c.Params("id")

	// AMAN: Konversi interface ke float64 dulu (bawaan JWT), baru ke uint
	userLocals := c.Locals("user_id")
	var userId uint
	if v, ok := userLocals.(float64); ok {
		userId = uint(v)
	} else if v, ok := userLocals.(uint); ok {
		userId = v
	} else {
		return c.Status(401).JSON(fiber.Map{"message": "User ID corrupt"})
	}

	// 1. Cek Toko User
	var toko models.Toko
	if err := database.DB.Where("user_id = ?", userId).First(&toko).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Toko tidak ditemukan / Belum bikin toko"})
	}

	// 2. Cek Barang
	var product models.Product
	if err := database.DB.First(&product, productId).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Barang gak ketemu"})
	}

	// 3. SECURITY CHECK: IDOR Protection
	if product.TokoID != toko.ID {
		return c.Status(403).JSON(fiber.Map{"message": "Jangan hapus dagangan orang woy!"})
	}

	// 4. Hapus (Soft Delete)
	database.DB.Delete(&product)

	return c.JSON(fiber.Map{"message": "Barang berhasil dihapus"})
}

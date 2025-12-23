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

// Struct Input buat nangkep JSON dari Postman
type CreateProductInput struct {
	NamaProduk    string  `json:"nama_produk"`
	HargaReseller float64 `json:"harga_reseller"`
	HargaKonsumen float64 `json:"harga_konsumen"`
	Stok          int     `json:"stok"`
	Deskripsi     string  `json:"deskripsi"`
	PhotoURL      string  `json:"photo_url"` // URL dari Upload Controller tadi
}

// CreateProduct - Logic Upload Barang
func CreateProduct(c *fiber.Ctx) error {
	// 1. Ambil User ID dari Token (Cara Aman)
	userLocals := c.Locals("user_id")
	var userId uint
	if v, ok := userLocals.(float64); ok {
		userId = uint(v)
	} else if v, ok := userLocals.(uint); ok { // Jaga-jaga kalau formatnya beda
		userId = v
	} else {
		return c.Status(401).JSON(fiber.Map{"message": "Unauthorized"})
	}

	// 2. Cari Toko punya User
	var toko models.Toko // Pake struct 'Toko'
	if err := database.DB.Where("user_id = ?", userId).First(&toko).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Lo belum punya toko bro!"})
	}

	// 3. Parsing Input
	var input CreateProductInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Input error, cek tipe data harga (jangan pake string)"})
	}

	// 4. Masukin ke Model Product
	newProduct := models.Product{
		NamaProduk:    input.NamaProduk,
		HargaReseller: input.HargaReseller,
		HargaKonsumen: input.HargaKonsumen,
		Stok:          input.Stok,
		Deskripsi:     input.Deskripsi,
		TokoID:        toko.ID, // Sambungin ke Toko
		Slug:          strings.ReplaceAll(strings.ToLower(input.NamaProduk), " ", "-") + "-" + strconv.Itoa(int(time.Now().Unix())),
	}

	// 5. Simpan Produk
	if err := database.DB.Create(&newProduct).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal simpan produk"})
	}

	// 6.  LOGIC PENTING: Simpan URL Foto ke Tabel FotoProduk 
	if input.PhotoURL != "" {
		foto := models.FotoProduk{
			ProductID: newProduct.ID,
			Url:       input.PhotoURL,
		}
		database.DB.Create(&foto)
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Barang berhasil dijual!",
		"data":    newProduct,
	})
}

// GetAllProducts - Katalog
func GetAllProducts(c *fiber.Ctx) error {
	var products []models.Product

	// Preload Toko & FotoProduk (Category jangan dulu)
	query := database.DB.Preload("Toko").Preload("FotoProduk")

	// Search Logic
	search := c.Query("search")
	if search != "" {
		query = query.Where("nama_produk LIKE ?", "%"+search+"%")
	}

	query.Find(&products)

	return c.JSON(fiber.Map{"data": products})
}

// GetProductDetail - Detail 1 Barang
func GetProductDetail(c *fiber.Ctx) error {
	productId := c.Params("id")
	var product models.Product

	if err := database.DB.Preload("Toko").Preload("FotoProduk").First(&product, productId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{"message": "Barang gak ketemu"})
		}
		return c.Status(500).JSON(fiber.Map{"message": "Error database"})
	}

	return c.JSON(fiber.Map{"data": product})
}

// UpdateProduct - Update Barang (Versi Aman + Support Ganti Harga)
func UpdateProduct(c *fiber.Ctx) error {
	productId := c.Params("id")

	// 1. Ambil User ID (Safety Check)
	userLocals := c.Locals("user_id")
	var userId uint
	if v, ok := userLocals.(float64); ok {
		userId = uint(v)
	} else if v, ok := userLocals.(uint); ok {
		userId = v
	} else {
		return c.Status(401).JSON(fiber.Map{"message": "User ID corrupt"})
	}

	// 2. Cek Toko User (Pake struct 'Toko')
	var toko models.Toko
	if err := database.DB.Where("user_id = ?", userId).First(&toko).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Toko tidak ditemukan"})
	}

	// 3. Cek Barang Existing
	var product models.Product
	if err := database.DB.First(&product, productId).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Barang gak ketemu"})
	}

	// 4. SECURITY CHECK: IDOR Protection
	// Pastiin barang ini beneran punya toko si user
	if product.TokoID != toko.ID {
		return c.Status(403).JSON(fiber.Map{"message": "Dih, mau ngedit barang orang ya? Gak boleh!"})
	}

	// 5. Nangkep Input Baru
	// Kita pake struct input yang sama kayak Create biar konsisten
	var input CreateProductInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Input error"})
	}

	// 6. Update Data
	// Kita update manual field-nya biar aman
	product.NamaProduk = input.NamaProduk
	product.HargaReseller = input.HargaReseller
	product.HargaKonsumen = input.HargaKonsumen
	product.Stok = input.Stok
	product.Deskripsi = input.Deskripsi
	// product.CategoryID = input.CategoryID (Nyalain kalau kategori dah ada)

	database.DB.Save(&product)

	// Note: Buat update foto agak kompleks (harus hapus lama -> insert baru).
	// Buat malem ini kita skip dulu update foto, fokus data teks aja.

	return c.JSON(fiber.Map{"message": "Barang berhasil diupdate", "data": product})
}

// DeleteProduct - Hapus Barang (Versi Aman)
func DeleteProduct(c *fiber.Ctx) error {
	productId := c.Params("id")

	// 1. Ambil User ID
	userLocals := c.Locals("user_id")
	var userId uint
	if v, ok := userLocals.(float64); ok {
		userId = uint(v)
	} else if v, ok := userLocals.(uint); ok {
		userId = v
	} else {
		return c.Status(401).JSON(fiber.Map{"message": "User ID corrupt"})
	}

	// 2. Cek Toko
	var toko models.Toko
	if err := database.DB.Where("user_id = ?", userId).First(&toko).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Toko tidak ditemukan"})
	}

	// 3. Cek Barang
	var product models.Product
	if err := database.DB.First(&product, productId).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Barang gak ketemu"})
	}

	// 4. SECURITY CHECK: IDOR
	if product.TokoID != toko.ID {
		return c.Status(403).JSON(fiber.Map{"message": "Jangan hapus dagangan orang woy!"})
	}

	// 5. Hapus (Soft Delete)
	database.DB.Delete(&product)

	return c.JSON(fiber.Map{"message": "Barang berhasil dihapus"})
}

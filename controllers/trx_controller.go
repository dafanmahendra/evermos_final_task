package controllers

import (
	"fmt"
	"time"

	"github.com/dafanmahendra/evermos-backend/database"
	"github.com/dafanmahendra/evermos-backend/models"
	"github.com/gofiber/fiber/v2"
)

// Struct buat nerima Input JSON dari Frontend
// User kirim: "Kirim ke Alamat ID 5, Barang yg dibeli ini list-nya..."
type CheckoutInput struct {
	AlamatID uint `json:"alamat_id"`
	Items    []struct {
		ProductID uint `json:"product_id"`
		Quantity  int  `json:"quantity"`
	} `json:"items"`
}

func Checkout(c *fiber.Ctx) error {
	// 1. Ambil User ID dari Token
	userLocals := c.Locals("user_id")
	var userId uint
	if v, ok := userLocals.(float64); ok {
		userId = uint(v)
	} else {
		return c.Status(401).JSON(fiber.Map{"message": "Unauthorized"})
	}

	// 2. Parsing Input JSON
	var input CheckoutInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Format order salah"})
	}

	// 3. MULAI TRANSAKSI DATABASE (Jurus Anti Error)
	// Kita pake 'tx' (Transaction), bukan 'database.DB' biasa.
	tx := database.DB.Begin()

	// Safety Net: Kalau ada panic/error parah, otomatis Rollback (Batalin semua)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var totalBelanja float64
	var detailTrxList []models.DetailTrx

	// 4. LOOPING BARANG (Proses satu per satu)
	for _, item := range input.Items {
		var product models.Product
		// Cek stok barang (Pake 'tx' buat query)
		if err := tx.First(&product, item.ProductID).Error; err != nil {
			tx.Rollback() // Batalin semua
			return c.Status(404).JSON(fiber.Map{"message": "Barang id " + fmt.Sprint(item.ProductID) + " ga ketemu"})
		}

		// Validasi Stok
		if product.Stok < item.Quantity {
			tx.Rollback() // Batalin semua
			return c.Status(400).JSON(fiber.Map{"message": "Stok " + product.NamaProduk + " abis bro!"})
		}

		//  FOTO COPY DATA (SNAPSHOT) ke LogProduk
		// Biar kalau harga asli berubah, history belanja tetep aman
		logProduk := models.LogProduk{
			ProdukAsliID:  product.ID,
			NamaProduk:    product.NamaProduk,
			FotoProduk:    "", // Nanti lo bisa ambil dari relasi foto kalau mau detail
			HargaReseller: product.HargaReseller,
			HargaKonsumen: product.HargaKonsumen,
			Deskripsi:     product.Deskripsi,
			TokoID:        product.TokoID,
		}
		if err := tx.Create(&logProduk).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"message": "Gagal bikin log produk"})
		}

		// Hitung Subtotal
		subTotal := product.HargaKonsumen * float64(item.Quantity)
		totalBelanja += subTotal

		// Siapin Detail Transaksi
		detail := models.DetailTrx{
			LogProdukID: logProduk.ID, // Link ke Log, bukan ke Produk Asli
			Kuantitas:   item.Quantity,
			HargaTotal:  subTotal,
		}
		detailTrxList = append(detailTrxList, detail)

		//  POTONG STOK ASLI
		product.Stok = product.Stok - item.Quantity
		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"message": "Gagal potong stok"})
		}
	}

	// 5. BIKIN HEADER TRANSAKSI (Nota Utama)
	newTrx := models.Trx{
		UserID:      userId,
		AlamatID:    input.AlamatID,
		HargaTotal:  totalBelanja,
		Status:      "UNPAID", // Status awal
		KodeInvoice: fmt.Sprintf("INV-%d-%d", userId, time.Now().Unix()), // Contoh: INV-1-1766333
		DetailTrx:   detailTrxList, // Masukin anak-anaknya tadi
	}

	if err := tx.Create(&newTrx).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"message": "Gagal bikin transaksi"})
	}

	// 6. COMMIT (RESMIKAN SEMUA PERUBAHAN)
	// Kalau sampai sini gak ada error, baru simpan permanen ke DB.
	tx.Commit()

	return c.JSON(fiber.Map{
		"message": "Checkout Berhasil!",
		"data":    newTrx,
	})
}
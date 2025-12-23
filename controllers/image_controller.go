package controllers

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// UploadImage handle upload file tunggal
func UploadImage(c *fiber.Ctx) error {
	// 1. Ambil file dari request (Key: "image")
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Gagal upload, pastikan key-nya 'image'"})
	}

	// 2. Validasi File (Harus Gambar)
	// Ambil ekstensi file (misal .jpg)
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return c.Status(400).JSON(fiber.Map{"message": "Hanya boleh upload file JPG/PNG"})
	}

	// 3. Cek Ukuran (Limit 2MB)
	if file.Size > 2*1024*1024 {
		return c.Status(400).JSON(fiber.Map{"message": "File terlalu besar (Maks 2MB)"})
	}

	// 4. Bikin Nama Unik (Pake Timestamp)
	// Contoh: 17098822_kucing.jpg
	filename := fmt.Sprintf("%d%s", time.Now().Unix(), ext)

	// 5. Simpan ke folder public/uploads
	path := fmt.Sprintf("./public/uploads/%s", filename)
	if err := c.SaveFile(file, path); err != nil {
		fmt.Println(" ERROR ASLI NYA INI: ", err) // <--- Tambahin baris ini
		return c.Status(500).JSON(fiber.Map{"message": "Gagal menyimpan file ke server"})
	}

	// 6. Generate URL buat Frontend
	// Nanti ini dipake buat disimpen di tabel Produk/User
	imageUrl := fmt.Sprintf("http://localhost:8080/uploads/%s", filename)

	return c.JSON(fiber.Map{
		"message": "Upload berhasil",
		"url":     imageUrl,
	})
}

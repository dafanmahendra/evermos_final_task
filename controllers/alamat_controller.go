package controllers

import (
    "github.com/dafanmahendra/evermos-backend/database"
    "github.com/dafanmahendra/evermos-backend/models"
    "github.com/gofiber/fiber/v2"
)

// CreateAlamat - Nambahin alamat baru (Wajib Login)
func CreateAlamat(c *fiber.Ctx) error {
    // 1. Ambil ID User dari Token (Satpam)
    userLocals := c.Locals("user_id")
    if userLocals == nil {
        return c.Status(401).JSON(fiber.Map{"message": "Unauthorized"})
    }
    // Casting float64 (default JWT) ke uint
    userId := uint(userLocals.(float64))

    // 2. Parse Input JSON
    var input models.Alamat
    if err := c.BodyParser(&input); err != nil {
        return c.Status(400).JSON(fiber.Map{"message": "Input tidak valid"})
    }

    // 3. Paksa UserID sesuai Token (Biar gak bisa ngisi alamat buat orang lain - Ketentuan No. 12)
    input.UserID = userId

    // 4. Simpan ke Database
    if err := database.DB.Create(&input).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"message": "Gagal menyimpan alamat"})
    }

    return c.Status(201).JSON(fiber.Map{
        "message": "Alamat berhasil ditambahkan",
        "data":    input,
    })
}

// GetMyAlamats - Liat list alamat sendiri
func GetMyAlamats(c *fiber.Ctx) error {
    // 1. Ambil ID User
    userLocals := c.Locals("user_id")
    if userLocals == nil {
        return c.Status(401).JSON(fiber.Map{"message": "Unauthorized"})
    }
    userId := uint(userLocals.(float64))

    // 2. Cari semua alamat milik user ini
    var alamats []models.Alamat
    if err := database.DB.Where("user_id = ?", userId).Find(&alamats).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"message": "Gagal mengambil data alamat"})
    }

    return c.JSON(fiber.Map{
        "data": alamats,
    })
}
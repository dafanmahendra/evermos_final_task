package controllers

import (
	"github.com/dafanmahendra/evermos-backend/database"
	"github.com/dafanmahendra/evermos-backend/models"
	"github.com/gofiber/fiber/v2"
)

// CreateCategory - Cuma Admin yang boleh akses ini nanti
func CreateCategory(c *fiber.Ctx) error {
	var input models.Category

	// 1. Ambil Input
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Input error bro"})
	}

	// 2. Simpan ke DB
	if err := database.DB.Create(&input).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal bikin kategori"})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Kategori berhasil dibuat (Mantap Admin!)",
		"data":    input,
	})
}

// GetAllCategories - Ini boleh diakses semua orang (biar user bisa milih kategori)
func GetAllCategories(c *fiber.Ctx) error {
	var categories []models.Category
	database.DB.Find(&categories)
	return c.JSON(fiber.Map{"data": categories})
}

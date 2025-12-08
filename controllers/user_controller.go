package controllers

import (
	"fmt"
	"strings"

	"github.com/dafanmahendra/evermos-backend/database"
	"github.com/dafanmahendra/evermos-backend/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// Register - Daftar user baru + Auto Create Toko
func Register(c *fiber.Ctx) error {
	var input models.User

	// 1. Parse Input
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// 2. Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.KataSandi), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}
	input.KataSandi = string(hashedPassword)

	// 3. Simpan User ke Database
	if err := database.DB.Create(&input).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create user. Email/Phone might be duplicate.",
			"details": err.Error(),
		})
	}

	// --- 🚨 START SUNTIKAN LOGIC TUGAS NO. 2 (AUTO TOKO) 🚨 ---

	// Logic: Nama Toko = "Toko [Nama User]"
	namaToko := fmt.Sprintf("Toko %s", input.Nama)
	// Logic: URL = "evermos.com/toko-[nama-user-tanpa-spasi]"
	slug := strings.ReplaceAll(strings.ToLower(input.Nama), " ", "-")
	urlToko := fmt.Sprintf("evermos.com/toko-%s", slug)

	toko := models.Toko{
		IdUser:   input.ID, // Ambil ID user yang baru aja dibuat
		NamaToko: namaToko,
		UrlToko:  urlToko,
	}

	// Simpan Toko Otomatis
	if err := database.DB.Create(&toko).Error; err != nil {
		// Kalau gagal bikin toko, kita log aja errornya (User tetep kebuat)
		fmt.Println("Gagal membuat toko otomatis:", err)
	}
	// --- END SUNTIKAN ---

	// 4. Return Response Sukses
	input.KataSandi = "" // Hide password

	return c.Status(201).JSON(fiber.Map{
		"message": "User registered and Store created successfully",
		"data": fiber.Map{ // Bungkus rapi pake 'data' sesuai style API umum
			"user": input,
			"toko": toko, // Balikin data toko biar keliatan di response
		},
	})
}

// Login - Autentikasi user
func Login(c *fiber.Ctx) error {
	var input struct {
		Email     string `json:"email"`
		KataSandi string `json:"kata_sandi"`
	}

	// Parse body
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Cari user berdasarkan email
	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.KataSandi), []byte(input.KataSandi)); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	// Login sukses (bisa tambahin JWT token di sini nanti)
	user.KataSandi = "" // Jangan return password

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"user":    user,
	})
}

// GetAllUsers - Ambil semua user (Admin only, bisa ditambahin middleware nanti)
func GetAllUsers(c *fiber.Ctx) error {
	var users []models.User

	if err := database.DB.Find(&users).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}

	// Hide passwords
	for i := range users {
		users[i].KataSandi = ""
	}

	return c.JSON(fiber.Map{
		"users": users,
	})
}

// GetUserByID - Ambil user berdasarkan ID
func GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	if err := database.DB.Preload("Toko").First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	user.KataSandi = "" // Hide password

	return c.JSON(fiber.Map{
		"user": user,
	})
}

// UpdateUser - Update data user
func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	// Cek apakah user ada
	if err := database.DB.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Parse input update
	var input models.User
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Update fields (kecuali password, nanti bikin endpoint terpisah)
	if input.Nama != "" {
		user.Nama = input.Nama
	}
	if input.NoTelp != "" {
		user.NoTelp = input.NoTelp
	}
	if input.TanggalLahir != "" {
		user.TanggalLahir = input.TanggalLahir
	}
	if input.Pekerjaan != "" {
		user.Pekerjaan = input.Pekerjaan
	}
	if input.Email != "" {
		user.Email = input.Email
	}
	if input.IdProvinsi != "" {
		user.IdProvinsi = input.IdProvinsi
	}
	if input.IdKota != "" {
		user.IdKota = input.IdKota
	}

	// Save ke database
	if err := database.DB.Save(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	user.KataSandi = "" // Hide password

	return c.JSON(fiber.Map{
		"message": "User updated successfully",
		"user":    user,
	})
}

// DeleteUser - Hapus user
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	// Cek apakah user ada
	if err := database.DB.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Delete user (CASCADE bakal hapus toko juga karena ada constraint di model)
	if err := database.DB.Delete(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

package database

import (
	"fmt"
	"github.com/dafanmahendra/evermos-backend/models" // Import package models lo
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// Config koneksi (User: root, Pass: secret, DB: inventory_db)
	dsn := "root:secret@tcp(127.0.0.1:3306)/inventory_db?charset=utf8mb4&parseTime=True&loc=Local"
	
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Gagal konek ke Database!")
	}

	fmt.Println("✅ Sukses Konek ke MySQL via GORM!")

	// Auto Migrate: GORM bakal bikin tabel users & toko otomatis sesuai struct tadi
	err = DB.AutoMigrate(&models.User{}, &models.Toko{})
	if err != nil {
		fmt.Println("Gagal Migrasi Tabel:", err)
	} else {
		fmt.Println("✅ Sukses Migrasi Tabel User & Toko!")
	}
}
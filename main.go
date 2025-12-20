package main

import (
	"github.com/dafanmahendra/evermos-backend/database"
	"github.com/dafanmahendra/evermos-backend/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// 1. Konek Database & Migrasi Tabel
	database.ConnectDB()

	// 2. Setup Fiber (Framework)
	app := fiber.New()

	// 3. Middleware
	app.Use(logger.New()) // Log setiap request
	app.Use(cors.New())   // Enable CORS untuk frontend

	// 4. Test Route Sederhana
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Server Evermos Ready!",
		})
	})

	// 5. Setup Routes API
	routes.Setup(app) 

	// 6. Jalanin Server di port 8080
	app.Listen(":8080")
}

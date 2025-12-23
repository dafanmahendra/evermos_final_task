package main

import (
	"fmt"
	"log"

	"github.com/dafanmahendra/evermos-backend/database"
	"github.com/dafanmahendra/evermos-backend/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// 1. Konek Database
	database.ConnectDB()

	// 2. Setup Fiber
	app := fiber.New()

	// 3. Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// --- SETUP STATIC FILES
	fmt.Println(" MENGAKTIFKAN FOLDER GAMBAR...")
	app.Static("/uploads", "./public/uploads")
	// ------------------------------------------------

	// 4. Test Route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Server Ready!"})
	})

	// 5. Setup Routes API
	routes.Setup(app)

	// 6. Jalanin Server
	fmt.Println("SERVER SIAP DI PORT :8080")
	log.Fatal(app.Listen(":8080"))
}

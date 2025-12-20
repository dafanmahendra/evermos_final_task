package routes

import (
	"github.com/dafanmahendra/evermos-backend/controllers"
	"github.com/dafanmahendra/evermos-backend/middleware" // Hampir lupa njir 
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/api/v1")

	// Auth Routes (Gak perlu satpam)
	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)

	// --- AREA BERBAHAYA (User Routes) ---
	// Kita pasang Satpam "Protected" di sini
	userGroup := api.Group("/users", middleware.Protected)

	userGroup.Get("/", controllers.GetAllUsers) // Sekarang butuh token!
	userGroup.Get("/:id", controllers.GetUserByID)
	userGroup.Put("/:id", controllers.UpdateUser)
	userGroup.Delete("/:id", controllers.DeleteUser)

	// --- CATEGORY ROUTES ---
	// Semua rute category butuh Login (Protected)
	categoryGroup := api.Group("/category", middleware.Protected)

	// GET: User biasa BOLEH liat list kategori
	categoryGroup.Get("/", controllers.GetAllCategories)

	// POST: CUMA ADMIN yang boleh bikin (Satpam Lapis 2)
	// Perhatikan urutannya: Protected (Cek Token) -> AdminOnly (Cek Role) -> Controller
	categoryGroup.Post("/", middleware.AdminOnly, controllers.CreateCategory)

	// --- PRODUCT ROUTES ---
	productGroup := api.Group("/products")

	// Public: Semua orang boleh liat (gak perlu login)
	productGroup.Get("/", controllers.GetAllProducts)
	productGroup.Get("/:id", controllers.GetProductDetail)

	// Protected: Harus Login buat create/update/delete
	productGroup.Post("/", middleware.Protected, controllers.CreateProduct)
	productGroup.Put("/:id", middleware.Protected, controllers.UpdateProduct)
	productGroup.Delete("/:id", middleware.Protected, controllers.DeleteProduct)
}

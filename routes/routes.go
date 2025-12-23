package routes

import (
	"github.com/dafanmahendra/evermos-backend/controllers"
	"github.com/dafanmahendra/evermos-backend/middleware" // Hampir lupa njir
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/api/v1")

	// Auth Routes gak perlu security check
	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)

	// upload route lah intinya
	api.Post("/upload", middleware.Protected, controllers.UploadImage)

	// CRUD BUAT BESOK
	// PRODUCT ROUTES

	// 1. Jualan (Create) - HARUS LOGIN
	api.Post("/products", middleware.Protected, controllers.CreateProduct)

	// 2. Liat Katalog (Get All) - Boleh publik (biar orang bisa window shopping)
	api.Get("/products", controllers.GetAllProducts)

	// 3. Liat Detail 1 Barang (Get One)
	api.Get("/products/:id", controllers.GetProductDetail)
	
	// 4. Edit Dagangan (Update) - HARUS LOGIN
	api.Put("/products/:id", middleware.Protected, controllers.UpdateProduct)
	
	// 5. Hapus Dagangan (Delete) - HARUS LOGIN
	api.Delete("/products/:id", middleware.Protected, controllers.DeleteProduct)
	
	// TRANSAKSI ROUTES BARU
	api.Post("/checkout", middleware.Protected, controllers.Checkout)

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

	// ALAMAT ROUTES BARU
	alamatGroup := api.Group("/alamats", middleware.Protected)

	alamatGroup.Post("/", controllers.CreateAlamat) // buat alamat
	alamatGroup.Get("/", controllers.GetMyAlamats)  // liat alamat sendiri

}

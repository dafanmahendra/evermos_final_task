package routes

import (
	"github.com/dafanmahendra/evermos-backend/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Group API v1
	v1 := app.Group("/api/v1")

	// Auth endpoints
	v1.Post("/register", controllers.Register) // POST /api/v1/register
	v1.Post("/login", controllers.Login)       // POST /api/v1/login

	// User CRUD endpoints
	v1.Get("/users", controllers.GetAllUsers)       // GET /api/v1/users
	v1.Get("/users/:id", controllers.GetUserByID)   // GET /api/v1/users/:id
	v1.Put("/users/:id", controllers.UpdateUser)    // PUT /api/v1/users/:id
	v1.Delete("/users/:id", controllers.DeleteUser) // DELETE /api/v1/users/:id
}




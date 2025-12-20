package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Protected - Satpam Pintu Utama (Cek Token Valid/Gak)
func Protected(c *fiber.Ctx) error {

	fmt.Println("Request Masuk: ", c.Path())

	// 1. Ambil token dari Header "Authorization"
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{"message": "Unauthorized: Token gak ada bro"})
	}

	// 2. Format biasanya "Bearer <token>", kita butuh ambil tokennya aja
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	// 3. Parse & Validasi Token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Pastikan method enkripsinya bener (HMAC)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}
		// Kunci rahasia harus SAMA PERSIS kayak di Login tadi
		return []byte("rahasia_negara"), nil
	})

	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{"message": "Unauthorized: Token gak valid atau expired"})
	}

	// 4. Kalau Token Valid, Ambil data user dari dalamnya
	claims := token.Claims.(jwt.MapClaims)
	userId := claims["user_id"]
	isAdmin := claims["isAdmin"]

	// 5. Simpan data user ke "Context" (Biar bisa dipake di Controller nanti)
	c.Locals("user_id", userId)
	c.Locals("isAdmin", isAdmin)

	// 6. Lanjut masuk ke dalam
	return c.Next()
}

// AdminOnly - Satpam VIP (Cek User Admin/Bukan)
func AdminOnly(c *fiber.Ctx) error {
	// Ambil status isAdmin yang udah disimpen Satpam Protected tadi
	isAdmin := c.Locals("isAdmin")

	// Cek: Kalau nil (gak ada) atau false, tendang!
	if isAdmin == nil || isAdmin == false {
		return c.Status(403).JSON(fiber.Map{"message": "Forbidden: Khusus Admin woy!"})
	}

	return c.Next()
}

// controller/auth_controller.go
package controller

import (
	"crud-alumni/app/models"
	"crud-alumni/app/service"

	"github.com/gofiber/fiber/v2"
)

// Struct untuk AuthController, bergantung pada AuthService
type AuthController struct {
	authService service.AuthService
}

// Constructor untuk AuthController
func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Method Login yang sudah di-refactor
func (ac *AuthController) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Request body tidak valid",
		})
	}

	// Validasi input
	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username dan password harus diisi",
		})
	}

	// Panggil service untuk menjalankan logika bisnis
	response, err := ac.authService.Login(req)
	if err != nil {
		// Service akan mengembalikan error yang sudah sesuai
		// Kita hanya perlu menentukan status code berdasarkan isi error
		if err.Error() == "username atau password salah" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Login berhasil",
		"data":    response,
	})
}

// Method getProfile bisa tetap di sini, karena hanya membaca dari context
func (ac *AuthController) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int) // Sesuaikan tipe data dengan model User.ID Anda
	username := c.Locals("username").(string)
	role := c.Locals("role").(string)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Profile berhasil diambil",
		"data": fiber.Map{
			"user_id":  userID,
			"username": username,
			"role":     role,
		},
	})
}
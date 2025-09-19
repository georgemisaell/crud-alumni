// controller/auth_controller.go
package controller

import (
	"crud-alumni/app/models"
	"crud-alumni/app/service"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

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

	response, err := ac.authService.Login(req)
	if err != nil {
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

func (ac *AuthController) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
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
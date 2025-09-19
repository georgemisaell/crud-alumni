package main

import (
	"crud-alumni/app/repository"
	"crud-alumni/app/service"
	"crud-alumni/controller"
	"crud-alumni/database"
	route "crud-alumni/routes"
	"log"

	// "time"

	"github.com/gofiber/fiber/v2"
)

func main(){
	// Koneksi database
	database.ConnectDB()
	defer database.DB.Close()

	// Inisialisasi Fiber
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error)error{
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Repository
	alumniRepo := repository.NewAlumniRepository(database.DB)
	pekerjaanRepo := repository.NewPekerjaanRepository(database.DB)
	userRepo := repository.NewUserRepository(database.DB)
	
	// Service
	alumniService := service.NewAlumniService(alumniRepo)
	pekerjaanService := service.NewPekerjaanService(pekerjaanRepo)
	authService := service.NewAuthService(userRepo)

	// Controller
	alumniController := controller.NewAlumniController(alumniService)
	pekerjaanController := controller.NewPekerjaanController(pekerjaanService)
	authController := controller.NewAuthController(authService)

	// Routes
	route.SetupRoutes(app, alumniController, pekerjaanController, authController)

	//start server
	log.Println("Server is running on port 3000...")
	log.Fatal(app.Listen(":3000"))
}

package main

import (
	"crud-alumni/database"
	"crud-alumni/route"
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

	// Routes
	route.SetupRoutes(app)

	//start server
	log.Fatal(app.Listen(":3000"))
}

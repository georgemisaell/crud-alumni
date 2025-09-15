package route

import (
	"crud-alumni/controller"
	"crud-alumni/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, alumniController *controller.AlumniController, pekerjaanController *controller.PekerjaanController, authController *controller.AuthController) {
	api := app.Group("/api")	

	// Auth Routes
	auth := api.Group("/auth")
	auth.Post("/login", authController.Login) // Gunakan method dari authController
	// Tambahkan route untuk profile, dilindungi oleh middleware
	auth.Get("/profile", middleware.AuthRequired(), authController.GetProfile)

	// Alumni routes
	alumni := api.Group("/alumni", middleware.AuthRequired())
	alumni.Get("/", alumniController.GetAllAlumni)
	alumni.Get("/:id", alumniController.GetAlumniByID)
	// Akses Admin
	alumni.Post("/", middleware.AdminOnly(), alumniController.CreateAlumni)
	alumni.Put("/:id", middleware.AdminOnly() , alumniController.UpdateAlumni)
	alumni.Delete("/:id", middleware.AdminOnly(), alumniController.DeleteAlumni)

	alumniPekerjaan := api.Group("/alumni/pekerjaan")
	alumniPekerjaan.Get("/:tahun_lulus", alumniController.GetAlumniByYear)

	// Pekerjaan alumni
	pekerjaan := api.Group("/pekerjaan", middleware.AuthRequired())
	pekerjaan.Get("/", pekerjaanController.GetAllPekerjaan)
	pekerjaan.Get("/:id", pekerjaanController.GetPekerjaanByID)
	// Akses Admin
	pekerjaan.Post("/", middleware.AdminOnly(), pekerjaanController.CreatePekerjaan)
	pekerjaan.Put("/:id", middleware.AdminOnly(), pekerjaanController.UpdatePekerjaan)
	pekerjaan.Delete("/:id", middleware.AdminOnly(), pekerjaanController.DeletePekerjaan)

}
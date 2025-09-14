package route

import (
	"crud-alumni/controller"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, alumniController *controller.AlumniController, pekerjaanController *controller.PekerjaanController) {
	api := app.Group("/api")

	// Alumni routes
	alumni := api.Group("/alumni")
	alumni.Get("/", alumniController.GetAllAlumni)
	alumni.Get("/:id", alumniController.GetAlumniByID)
	alumni.Post("/", alumniController.CreateAlumni)
	alumni.Put("/:id", alumniController.UpdateAlumni)
	alumni.Delete("/:id", alumniController.DeleteAlumni)
	alumniPekerjaan := api.Group("/alumni/pekerjaan")
	alumniPekerjaan.Get("/:tahun_lulus", alumniController.GetAlumniByYear)

	// Pekerjaan alumni
	pekerjaan := api.Group("/pekerjaan")
	pekerjaan.Get("/", pekerjaanController.GetAllPekerjaan)
	pekerjaan.Get("/:id", pekerjaanController.GetPekerjaanByID)
	pekerjaan.Post("/", pekerjaanController.CreatePekerjaan)
	pekerjaan.Put("/:id", pekerjaanController.UpdatePekerjaan)
	pekerjaan.Delete("/:id", pekerjaanController.DeletePekerjaan)

}
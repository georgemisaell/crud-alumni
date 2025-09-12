package route

import (
	"crud-alumni/app/repository"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Alumni routes
	alumni := api.Group("/alumni")
	alumni.Get("/", repository.GetAllAlumni)
	alumni.Get("/:id", repository.GetAllAlumniByid)
	alumni.Post("/", repository.CreateAlumni)
	alumni.Put("/:id", repository.UpdateAlumni)
	alumni.Delete("/:id", repository.DeleteAlumni)

	// Pekerjaan alumni
	pekerjaan := api.Group("/pekerjaan")
	pekerjaan.Get("/", repository.GetAllPekerjaan)
	pekerjaan.Get("/:id", repository.GetAllPekerjaanByid)
	pekerjaan.Post("/", repository.CreatePekerjaanAlumni)
	pekerjaan.Put("/:id", repository.UpdatePekerjaanAlumni)
	pekerjaan.Delete("/:id", repository.DeletePekerjaanAlumni)

	alumniPekerjaan := api.Group("/alumni/pekerjaan")
	alumniPekerjaan.Get("/:tahun_lulus", repository.GetAlumniByYear)
}
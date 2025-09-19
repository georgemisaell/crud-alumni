package controller

import (
	"crud-alumni/app/models"
	"crud-alumni/app/service"
	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// AlumniController menangani semua interaksi HTTP.
type AlumniController struct {
	// Controller bergantung pada INTERFACE service.
	alumniService service.AlumniService
}

// NewAlumniController membuat instance baru dari AlumniController.
func NewAlumniController(service service.AlumniService) *AlumniController {
	return &AlumniController{alumniService: service}
}

func (ctrl *AlumniController) GetAllAlumni(c *fiber.Ctx) error {
	alumniList, err := ctrl.alumniService.GetAllAlumni()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data dari database",
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    alumniList,
		"message": "Data alumni berhasil diambil!",
	})
}

func (ctrl *AlumniController) GetAlumniByYear(c *fiber.Ctx) error {
	tahunLulus, err := strconv.Atoi(c.Params("tahun_lulus"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Parameter tahun lulus tidak valid",
		})
	}
	results, err := ctrl.alumniService.GetAlumniByYear(tahunLulus)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Terjadi kesalahan pada server saat mengambil data",
		})
	}
	if len(results) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "Data alumni tidak ditemukan untuk tahun kelulusan tersebut",
			"data":    []models.AlumniPekerjaan{},
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Data alumni berhasil diambil",
		"data":    results,
	})
}

func (ctrl *AlumniController) GetAlumniByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "ID tidak valid",
		})
	}

	alumni, err := ctrl.alumniService.GetAlumniByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false, "message": "Alumni tidak ditemukan",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Kesalahan server internal",
		})
	}
	return c.JSON(fiber.Map{
		"success": true, "data": alumni, "message": "Data alumni berhasil diambil",
	})
}

func (ctrl *AlumniController) CreateAlumni(c *fiber.Ctx) error {
	var req models.CreateAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Request body tidak valid",
		})
	}

	// Validasi input sederhana (apakah field kosong) ada di Controller.
	if req.NIM == "" || req.Nama == "" || req.Jurusan == "" || req.Email == "" || req.Angkatan == 0 || req.NoTelepon == "" || req.TahunLulus == 0 || req.Alamat == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Semua field harus diisi",
		})
	}

	newAlumni, err := ctrl.alumniService.CreateAlumni(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal menambah data alumni. Pastikan data unik (NIM/Email) belum digunakan.",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true, "data": newAlumni, "message": "Mahasiswa berhasil ditambahkan",
	})
}

func (ctrl *AlumniController) UpdateAlumni(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}
	var req models.UpdateAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Request body tidak valid"})
	}
	if req.NIM == "" || req.Nama == "" { // contoh validasi singkat
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Field wajib harus diisi"})
	}
	updatedAlumni, err := ctrl.alumniService.UpdateAlumni(id, req)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "Alumni tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Gagal mengupdate data"})
	}
	return c.JSON(fiber.Map{"success": true, "data": updatedAlumni, "message": "Alumni berhasil diupdate"})
}

func (ctrl *AlumniController) DeleteAlumni(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}
	err = ctrl.alumniService.DeleteAlumni(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "Alumni tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Gagal menghapus alumni"})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Alumni berhasil dihapus"})
}

func (ctrl *AlumniController) GetAlumniController(c *fiber.Ctx) error {
	// 1. Ambil query parameter dari request
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "id")
	order := c.Query("order", "asc")
	search := c.Query("search", "")

	// 2. Panggil service (tidak ada lagi logika bisnis di sini)
	alumni, total, err := ctrl.alumniService.GetAlumniWithPagination(search, sortBy, order, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data alumni",
			"error":   err.Error(),
		})
	}

	// 3. Bangun response
	meta := models.MetaInfo{
		Page:   page,
		Limit:  limit,
		Total:  total,
		Pages:  (total + limit - 1) / limit, // Kalkulasi total halaman
		SortBy: sortBy,
		Order:  order,
		Search: search,
	}

	response := models.AlumniResponse{
		Data: alumni,
		Meta: meta,
	}

	return c.JSON(response)
}
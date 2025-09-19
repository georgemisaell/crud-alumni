package controller

import (
	"crud-alumni/app/models"
	"crud-alumni/app/service"
	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PekerjaanController struct {
	pekerjaanService service.PekerjaanService
}

func NewPekerjaanController(service service.PekerjaanService) *PekerjaanController {
	return &PekerjaanController{pekerjaanService: service}
}

func (ctrl *PekerjaanController) GetAllPekerjaan(c *fiber.Ctx) error {
	pekerjaanList, err := ctrl.pekerjaanService.GetAllPekerjaan()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data pekerjaan dari database",
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    pekerjaanList,
		"message": "Data pekerjaan berhasil diambil",
	})
}

func (ctrl *PekerjaanController) GetPekerjaanByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "ID tidak valid",
		})
	}

	pekerjaan, err := ctrl.pekerjaanService.GetPekerjaanByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false, "message": "Data pekerjaan tidak ditemukan",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Kesalahan server internal",
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    pekerjaan,
		"message": "Data pekerjaan berhasil diambil",
	})
}

func (ctrl *PekerjaanController) CreatePekerjaan(c *fiber.Ctx) error {
	var req models.CreatePekerjaanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Request body tidak valid",
		})
	}
	if req.AlumniID == 0 || req.NamaPerusahaan == "" || req.PosisiJabatan == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Field wajib (AlumniID, NamaPerusahaan, PosisiJabatan) harus diisi",
		})
	}

	newPekerjaan, err := ctrl.pekerjaanService.CreatePekerjaan(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal menambah data pekerjaan", "error": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    newPekerjaan,
		"message": "Pekerjaan berhasil ditambahkan",
	})
}

func (ctrl *PekerjaanController) UpdatePekerjaan(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "ID tidak valid",
		})
	}
	var req models.UpdatePekerjaanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Request body tidak valid",
		})
	}
	if req.AlumniID == 0 || req.NamaPerusahaan == "" || req.PosisiJabatan == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Field wajib (AlumniID, NamaPerusahaan, PosisiJabatan) harus diisi",
		})
	}

	updatedPekerjaan, err := ctrl.pekerjaanService.UpdatePekerjaan(id, req)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false, "message": "Pekerjaan tidak ditemukan untuk diupdate",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal mengupdate data pekerjaan", "error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    updatedPekerjaan,
		"message": "Pekerjaan berhasil diupdate",
	})
}

func (ctrl *PekerjaanController) DeletePekerjaan(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "ID tidak valid",
		})
	}

	err = ctrl.pekerjaanService.DeletePekerjaan(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false, "message": "Pekerjaan tidak ditemukan untuk dihapus",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal menghapus pekerjaan",
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Pekerjaan berhasil dihapus",
	})
}

func (ctrl *PekerjaanController) GetPekerjaanController(c *fiber.Ctx) error {
	// 1. Ambil query parameter dari request
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "id")
	order := c.Query("order", "asc")
	search := c.Query("search", "")

	pekerjaan, total, err := ctrl.pekerjaanService.GetWithPagination(search, sortBy, order, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data pekerjaan",
			"error":   err.Error(),
		})
	}

	meta := models.MetaInfo{
		Page:   page,
		Limit:  limit,
		Total:  total,
		Pages:  (total + limit - 1) / limit,
		SortBy: sortBy,
		Order:  order,
		Search: search,
	}

	response := models.PekerjaanResponse{
		Data: pekerjaan,
		Meta: meta,
	}

	return c.JSON(response)
}
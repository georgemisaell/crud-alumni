package repository

import (
	"crud-alumni/app/models"
	"crud-alumni/database"
	"time"

	"strconv"
	//"time"

	"github.com/gofiber/fiber/v2"
)

func GetAllPekerjaan(c *fiber.Ctx) error{
	rows, err := database.DB.Query(`SELECT * FROM pekerjaan_alumni`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": string(err.Error()),
		})
	}
	defer rows.Close()

	var PekerjaanAlumniList []models.PekerjaanAlumni

	for rows.Next(){
		var p models.PekerjaanAlumni
		err := rows.Scan(
			&p.ID, 
			&p.AlumniID,
			&p.NamaPerusahaan,
			&p.PosisiJabatan,
			&p.BidangIndustri,
			&p.LokasiKerja,
			&p.GajiRange,
			&p.TanggalMulaiKerja,
			&p.TanggalSelesaiKerja,
			&p.StatusPekerjaan,
			&p.DeskripsiPekerjaan,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil{
			return c.Status(500).JSON(fiber.Map{
				"error": "Gagal scan data alumni",
			})
		}
		PekerjaanAlumniList = append(PekerjaanAlumniList, p)
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data": PekerjaanAlumniList,
		"message" : "Data alumni berhasil diambil!",
	})
}

func GetAllPekerjaanByid(c *fiber.Ctx)error{
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error" : "ID tidak valid",
		})
	}

	var p models.PekerjaanAlumni
	row := database.DB.QueryRow(`
		SELECT * 
		FROM pekerjaan_alumni
		WHERE id = $1
	`,id)

	err = row.Scan(
		&p.ID, 
		&p.AlumniID,
		&p.NamaPerusahaan,
		&p.PosisiJabatan,
		&p.BidangIndustri,
		&p.LokasiKerja,
		&p.GajiRange,
		&p.TanggalMulaiKerja,
		&p.TanggalSelesaiKerja,
		&p.StatusPekerjaan,
		&p.DeskripsiPekerjaan,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Alumni tidak ditemukan",
		})
	}

	return  c.JSON(fiber.Map{
		"success": true,
		"data": p,
		"message" : "Data alumni berhasil diambil",
	})
}

func CreatePekerjaanAlumni(c *fiber.Ctx)error{
	var req models.CreatePekerjaanRequest

	if err := c.BodyParser(&req); err != nil{
		return c.Status(400).JSON(fiber.Map{
			"error": "Request body tidak valid",
		})
	}

	// validasi input
	if req.AlumniID == 0 || req.NamaPerusahaan == "" || req.PosisiJabatan == "" || req.BidangIndustri == "" || req.LokasiKerja == "" || req.TanggalMulaiKerja == ""{
		return c.Status(400).JSON(fiber.Map{
			"error": "Semua field harus diisi",
		})
	}
 
	var id int
	err := database.DB.QueryRow(`
		INSERT INTO pekerjaan_alumni (alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, tanggal_mulai_kerja, gaji_range, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`, req.AlumniID, req.NamaPerusahaan, req.PosisiJabatan, req.BidangIndustri, req.LokasiKerja, req.TanggalMulaiKerja, req.GajiRange, req.TanggalSelesaiKerja, req.StatusPekerjaan, req.DeskripsiPekerjaan, time.Now(), time.Now(),
	).Scan(&id)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error" : string(err.Error()),
		})
	}

	// Ambil data yang baru ditambahakn 
	var newPekerjaanAlumni models.PekerjaanAlumni
	row := database.DB.QueryRow(`
		SELECT alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, tanggal_mulai_kerja, gaji_range, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at
		FROM pekerjaan_alumni
		WHERE id = $1
	`, id)

	row.Scan(
		&newPekerjaanAlumni.AlumniID, 
		&newPekerjaanAlumni.NamaPerusahaan, 
		&newPekerjaanAlumni.PosisiJabatan,
		&newPekerjaanAlumni.BidangIndustri,
		&newPekerjaanAlumni.LokasiKerja,
		&newPekerjaanAlumni.TanggalMulaiKerja,
		&newPekerjaanAlumni.GajiRange,
		&newPekerjaanAlumni.TanggalSelesaiKerja,
		&newPekerjaanAlumni.StatusPekerjaan,
		&newPekerjaanAlumni.DeskripsiPekerjaan,
		&newPekerjaanAlumni.CreatedAt,
		&newPekerjaanAlumni.UpdatedAt,
	)
	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"data": newPekerjaanAlumni,
		"message": "Pekerjaan berhasil ditambahkan",
	})
}

func UpdatePekerjaanAlumni (c *fiber.Ctx)error{
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error" : "ID tidak valid",
		})
	}

	var req models.UpdatePekerjaanRequest
	if err := c.BodyParser(&req); err != nil{
		return c.Status(400).JSON(fiber.Map{
			"error": string(err.Error()),
		})
	}

	// validasi input
	if req.AlumniID == 0 || req.NamaPerusahaan == "" || req.PosisiJabatan == "" || req.BidangIndustri == "" || req.LokasiKerja == "" || req.TanggalMulaiKerja == ""{
		return c.Status(400).JSON(fiber.Map{
			"error": "Semua field harus diisi",
		})
	}

	result, err := database.DB.Exec(`
		UPDATE pekerjaan_alumni
		SET alumni_id = $1, nama_perusahaan = $2, posisi_jabatan = $3, bidang_industri = $4, lokasi_kerja = $5, tanggal_mulai_kerja = $6, gaji_range = $7, tanggal_selesai_kerja = $8, status_pekerjaan = $9, deskripsi_pekerjaan = $10, updated_at = $11
		WHERE id = $12
	`, req.AlumniID, req.NamaPerusahaan, req.PosisiJabatan, req.BidangIndustri, req.LokasiKerja, req.TanggalMulaiKerja, req.GajiRange, req.TanggalSelesaiKerja, req.StatusPekerjaan, req.DeskripsiPekerjaan, time.Now(), id)
	
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":  "yo",
		})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "Pekerjaan tidak ditemukan",
		})
	}

	// ambil data yang sudah di update
	var updatedPekerjaan models.PekerjaanAlumni
	row := database.DB.QueryRow(`
		SELECT *
		FROM pekerjaan_alumni
		WHERE id = $1
	`,id)

	row.Scan(
		&updatedPekerjaan.AlumniID, 
		&updatedPekerjaan.NamaPerusahaan, 
		&updatedPekerjaan.PosisiJabatan,
		&updatedPekerjaan.BidangIndustri,
		&updatedPekerjaan.LokasiKerja,
		&updatedPekerjaan.TanggalMulaiKerja,
		&updatedPekerjaan.GajiRange,
		&updatedPekerjaan.TanggalSelesaiKerja,
		&updatedPekerjaan.StatusPekerjaan,
		&updatedPekerjaan.DeskripsiPekerjaan,
		&updatedPekerjaan.CreatedAt,
		&updatedPekerjaan.UpdatedAt,
	)

	return c.JSON(fiber.Map{
		"success": true,
		"data": updatedPekerjaan,
		"message": "Alumni berhasil di update",
	})
}

func DeletePekerjaanAlumni (c *fiber.Ctx)error{
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil{
		return c.Status(400).JSON(fiber.Map{
			"error": "ID tidak valid",
		})
	}

	result, err := database.DB.Exec("DELETE FROM pekerjaan_alumni WHERE id = $1", id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
		"error": "Gagal menghapus pekerjaan",
		})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
		"error": "Alumni tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Pekerjaan berhasil dihapus",
	})

}
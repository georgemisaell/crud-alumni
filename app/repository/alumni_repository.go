package repository

import (
	"crud-alumni/app/models"
	"crud-alumni/database"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetAllAlumni(c *fiber.Ctx) error{
	rows, err := database.DB.Query(`SELECT * FROM alumni`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": string(err.Error()),
		})
	}
	defer rows.Close()

	var alumniList []models.Alumni

	for rows.Next(){
		var a models.Alumni
		err := rows.Scan(
			&a.ID, 
			&a.NIM, 
			&a.Nama,
			&a.Jurusan,
			&a.Angkatan,
			&a.TahunLulus,
			&a.Email,
			&a.NoTelepon,
			&a.Alamat,
			&a.CreatedAt,
			&a.UpdatedAt,
		)
		if err != nil{
			return c.Status(500).JSON(fiber.Map{
				"error": "Gagal scan data alumni",
			})
		}
		alumniList = append(alumniList, a)
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data": alumniList,
		"message" : "Data alumni berhasil diambil!",
	})
}

func GetAlumniByYear(c *fiber.Ctx) error {
	tahunLulus, err := strconv.Atoi(c.Params("tahun_lulus"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Parameter tahun lulus tidak valid",
		})
	}

	query := `
        SELECT 
            a.id, a.nama, a.jurusan, a.tahun_lulus, 
            pa.bidang_industri, pa.nama_perusahaan, pa.posisi_jabatan, pa.gaji_range,
            COUNT(a.id) OVER (PARTITION BY a.tahun_lulus) AS jumlah_alumni_per_tahun_lulus
        FROM alumni a
        INNER JOIN pekerjaan_alumni pa ON a.id = pa.alumni_id
        WHERE 
            a.tahun_lulus = $1
            AND CAST(
                split_part(
                    regexp_replace(pa.gaji_range, '[Rp. ]', '', 'g'),
                '-', 1)
            AS BIGINT) > 4000000
        ORDER BY a.nama ASC;
    `
	rows, err := database.DB.Query(query, tahunLulus)
	if err != nil {
		log.Printf("Error executing query: %v\n", err) 
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Terjadi kesalahan pada server",
		})
	}
	defer rows.Close()

	
	var results []models.AlumniPekerjaan
	
	for rows.Next() {
		var a models.AlumniPekerjaan 

		err := rows.Scan(
			&a.ID,
			&a.Nama,
			&a.Jurusan,
			&a.TahunLulus,
			&a.BidangIndustri,
			&a.NamaPerusahaan,
			&a.PosisiJabatan,
			&a.GajiRange,
			&a.JumlahAlumniPerTahun, 
		)

		if err != nil {
			log.Printf("Error scanning row: %v\n", err) // Log error untuk debug
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Gagal memproses data dari database",
			})
		}

		results = append(results, a)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error during rows iteration: %v\n", err) // Log error untuk debug
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal membaca hasil dari database",
		})
	}
	
	if len(results) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message":   "Data alumni tidak ditemukan untuk tahun kelulusan tersebut",
			"data":    []models.AlumniPekerjaan{}, 
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Data alumni berhasil diambil",
		"data":    results, 
	})
}

func GetAllAlumniByid(c *fiber.Ctx)error{
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error" : "ID tidak valid",
		})
	}

	var a models.Alumni
	row := database.DB.QueryRow(`
		SELECT * 
		FROM alumni
		WHERE id = $1
	`,id)

	err = row.Scan(
		&a.ID, 
		&a.NIM, 
		&a.Nama,
		&a.Jurusan,
		&a.Angkatan,
		&a.TahunLulus,
		&a.Email,
		&a.NoTelepon,
		&a.Alamat,
		&a.CreatedAt,
		&a.UpdatedAt,
	)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Alumni tidak ditemukan",
		})
	}

	return  c.JSON(fiber.Map{
		"success": true,
		"data": a,
		"message" : "Data alumni berhasil diambil",
	})
}

func CreateAlumni(c *fiber.Ctx)error{
	var req models.CreateAlumniRequest

	if err := c.BodyParser(&req); err != nil{
		return c.Status(400).JSON(fiber.Map{
			"error": "Request body tidak valid",
		})
	}

	// validasi input
	if req.NIM == "" || req.Nama == "" || req.Jurusan == "" || req.Email == "" || req.Angkatan == 0 || req.NoTelepon == ""|| req.TahunLulus == 0 || req.Alamat == ""{
		return c.Status(400).JSON(fiber.Map{
			"error": "Semua field harus diisi",
		})
	}
 
	var id int
	err := database.DB.QueryRow(`
		INSERT INTO alumni (nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`, req.NIM, req.Nama, req.Jurusan, req.Angkatan, req.TahunLulus, req.Email, req.NoTelepon, req.Alamat, time.Now(), time.Now(),
	).Scan(&id)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error" : "Gagal menambah data alumni. Pastikan data belum digunakan",
		})
	}

	// Ambil data yang baru ditambahakn 
	var newAlumni models.Alumni
	row := database.DB.QueryRow(`
		SELECT nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at
		FROM alumni
		WHERE id = $1
	`, id)

	row.Scan(
		&newAlumni.ID, 
		&newAlumni.NIM, 
		&newAlumni.Nama,
		&newAlumni.Jurusan,
		&newAlumni.Angkatan,
		&newAlumni.TahunLulus,
		&newAlumni.Email,
		&newAlumni.NoTelepon,
		&newAlumni.Alamat,
		&newAlumni.CreatedAt,
		&newAlumni.UpdatedAt,
	)
	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"data": newAlumni,
		"message": "Mahasiswa berhasil ditambahkan",
	})
}

func UpdateAlumni (c *fiber.Ctx)error{
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error" : "ID tidak valid",
		})
	}

	var req models.UpdateAlumniRequest
	if err := c.BodyParser(&req); err != nil{
		return c.Status(400).JSON(fiber.Map{
			"error": string(err.Error()),
		})
	}

	// validasi input
	if req.NIM == "" || req.Nama == "" || req.Jurusan == "" || req.Email == "" || req.Angkatan == 0 || req.NoTelepon == ""|| req.TahunLulus == 0 || req.Alamat == ""{
		return c.Status(400).JSON(fiber.Map{
			"error": "Semua field harus diisi",
		})
	}

	result, err := database.DB.Exec(`
		UPDATE alumni
		SET nim = $1, nama = $2, jurusan = $3, angkatan = $4, tahun_lulus = $5, email = $6, no_telepon = $7, alamat = $8, updated_at = $9
		WHERE id = $10
	`, req.NIM, req.Nama, req.Jurusan, req.Angkatan, req.TahunLulus, req.Email, req.NoTelepon, req.Alamat, time.Now(), id)
	
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":  string(err.Error()),
		})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "Alumni tidak ditemukan",
		})
	}

	// ambil data yang sudah di update
	var updatedAlumni models.Alumni
	row := database.DB.QueryRow(`
		SELECT *
		FROM alumni
		WHERE id = $1
	`,id)

	row.Scan(
		&updatedAlumni.ID, 
		&updatedAlumni.NIM, 
		&updatedAlumni.Nama,
		&updatedAlumni.Jurusan,
		&updatedAlumni.Angkatan,
		&updatedAlumni.TahunLulus,
		&updatedAlumni.Email,
		&updatedAlumni.NoTelepon,
		&updatedAlumni.Alamat,
		&updatedAlumni.CreatedAt,
		&updatedAlumni.UpdatedAt,
	)

	return c.JSON(fiber.Map{
		"success": true,
		"data": updatedAlumni,
		"message": "Alumni berhasil di update",
	})
}

func DeleteAlumni (c *fiber.Ctx)error{
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil{
		return c.Status(400).JSON(fiber.Map{
			"error": "ID tidak valid",
		})
	}

	result, err := database.DB.Exec("DELETE FROM alumni WHERE id = $1", id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
		"error": "Gagal menghapus alumni",
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
		"message": "Alumni berhasil dihapus",
	})

}
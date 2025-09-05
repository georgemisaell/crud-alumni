package main

import (
	"crud-alumni/app/models"
	"crud-alumni/database"
	"log"
	"strconv"
	"time"

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
	setupRoutes(app)

	//start server
	log.Fatal(app.Listen(":3000"))
}

func setupRoutes(app *fiber.App){
	api := app.Group("/api")

	// Alumni routes
	alumni := api.Group("/alumni")
	alumni.Get("/", getAllAlumni)
	alumni.Get("/:id", getAllAlumniByid)
	alumni.Post("/", CreateAlumni)

	// Pekerjaan alumni
}

func getAllAlumni(c *fiber.Ctx) error{
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

func getAllAlumniByid(c *fiber.Ctx)error{
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
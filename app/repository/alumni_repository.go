package repository

import (
	"crud-alumni/app/models"
	"database/sql"
	"log"
	"time"
)

// AlumniRepository mendefinisikan interface untuk operasi database alumni.
// Penggunaan interface memudahkan untuk testing (mocking).
type AlumniRepository interface {
	FindAll() ([]models.Alumni, error)
	FindByYear(year int) ([]models.AlumniPekerjaan, error)
	FindByID(id int) (*models.Alumni, error)
	Create(req models.CreateAlumniRequest) (*models.Alumni, error)
	Update(id int, req models.UpdateAlumniRequest) (*models.Alumni, error)
	Delete(id int) error
}

// alumniRepository adalah implementasi konkret dari AlumniRepository.
type alumniRepository struct {
	db *sql.DB
}

// NewAlumniRepository membuat instance baru dari alumniRepository.
func NewAlumniRepository(db *sql.DB) AlumniRepository {
	return &alumniRepository{db: db}
}

func (r *alumniRepository) FindAll() ([]models.Alumni, error) {
	rows, err := r.db.Query(`SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at FROM alumni`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alumniList []models.Alumni
	for rows.Next() {
		var a models.Alumni
		err := rows.Scan(
			&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan,
			&a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat,
			&a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return nil, err // Kembalikan error jika scan gagal
		}
		alumniList = append(alumniList, a)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return alumniList, nil
}

func (r *alumniRepository) FindByYear(year int) ([]models.AlumniPekerjaan, error) {
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
	rows, err := r.db.Query(query, year)
	if err != nil {
		log.Printf("Error executing query: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var results []models.AlumniPekerjaan
	for rows.Next() {
		var a models.AlumniPekerjaan
		err := rows.Scan(
			&a.ID, &a.Nama, &a.Jurusan, &a.TahunLulus, &a.BidangIndustri,
			&a.NamaPerusahaan, &a.PosisiJabatan, &a.GajiRange, &a.JumlahAlumniPerTahun,
		)
		if err != nil {
			log.Printf("Error scanning row: %v\n", err)
			return nil, err
		}
		results = append(results, a)
	}
	
	if err = rows.Err(); err != nil {
		log.Printf("Error during rows iteration: %v\n", err)
		return nil, err
	}

	return results, nil
}

func (r *alumniRepository) FindByID(id int) (*models.Alumni, error) {
	var a models.Alumni
	row := r.db.QueryRow(`
        SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at 
        FROM alumni WHERE id = $1
    `, id)

	err := row.Scan(
		&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan,
		&a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat,
		&a.CreatedAt, &a.UpdatedAt,
	)

	if err != nil {
		// Jika tidak ada baris yang ditemukan, sql.ErrNoRows akan dikembalikan
		return nil, err
	}

	return &a, nil
}

func (r *alumniRepository) Create(req models.CreateAlumniRequest) (*models.Alumni, error) {
	var newAlumni models.Alumni
	query := `
        INSERT INTO alumni (nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at
    `
	err := r.db.QueryRow(
		query,
		req.NIM, req.Nama, req.Jurusan, req.Angkatan, req.TahunLulus,
		req.Email, req.NoTelepon, req.Alamat, time.Now(), time.Now(),
	).Scan(
		&newAlumni.ID, &newAlumni.NIM, &newAlumni.Nama, &newAlumni.Jurusan, &newAlumni.Angkatan,
		&newAlumni.TahunLulus, &newAlumni.Email, &newAlumni.NoTelepon, &newAlumni.Alamat,
		&newAlumni.CreatedAt, &newAlumni.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &newAlumni, nil
}

func (r *alumniRepository) Update(id int, req models.UpdateAlumniRequest) (*models.Alumni, error) {
	var updatedAlumni models.Alumni
	query := `
        UPDATE alumni
        SET nim = $1, nama = $2, jurusan = $3, angkatan = $4, tahun_lulus = $5, email = $6, no_telepon = $7, alamat = $8, updated_at = $9
        WHERE id = $10
        RETURNING id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at
    `
	err := r.db.QueryRow(
		query,
		req.NIM, req.Nama, req.Jurusan, req.Angkatan, req.TahunLulus,
		req.Email, req.NoTelepon, req.Alamat, time.Now(), id,
	).Scan(
		&updatedAlumni.ID, &updatedAlumni.NIM, &updatedAlumni.Nama, &updatedAlumni.Jurusan, &updatedAlumni.Angkatan,
		&updatedAlumni.TahunLulus, &updatedAlumni.Email, &updatedAlumni.NoTelepon, &updatedAlumni.Alamat,
		&updatedAlumni.CreatedAt, &updatedAlumni.UpdatedAt,
	)

	if err != nil {
		// sql.ErrNoRows jika ID tidak ditemukan
		return nil, err
	}

	return &updatedAlumni, nil
}

func (r *alumniRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM alumni WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		// Mengembalikan error standar jika tidak ada baris yang terpengaruh
		return sql.ErrNoRows
	}

	return nil
}
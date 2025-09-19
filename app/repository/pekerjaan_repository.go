package repository

import (
	"crud-alumni/app/models"
	"database/sql"
	"time"
)

type PekerjaanRepository interface {
	FindAll() ([]models.PekerjaanAlumni, error)
	FindByID(id int) (*models.PekerjaanAlumni, error)
	Create(req models.CreatePekerjaanRequest) (*models.PekerjaanAlumni, error)
	Update(id int, req models.UpdatePekerjaanRequest) (*models.PekerjaanAlumni, error)
	Delete(id int) error
}

type pekerjaanRepository struct {
	db *sql.DB
}

func NewPekerjaanRepository(db *sql.DB) PekerjaanRepository {
	return &pekerjaanRepository{db: db}
}

func (r *pekerjaanRepository) FindAll() ([]models.PekerjaanAlumni, error) {
	rows, err := r.db.Query(`SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at FROM pekerjaan_alumni`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pekerjaanList []models.PekerjaanAlumni
	for rows.Next() {
		var p models.PekerjaanAlumni
		err := rows.Scan(
			&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri,
			&p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja,
			&p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		pekerjaanList = append(pekerjaanList, p)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pekerjaanList, nil
}

func (r *pekerjaanRepository) FindByID(id int) (*models.PekerjaanAlumni, error) {
	var p models.PekerjaanAlumni
	row := r.db.QueryRow(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at 
		FROM pekerjaan_alumni WHERE id = $1
	`, id)

	err := row.Scan(
		&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri,
		&p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja,
		&p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *pekerjaanRepository) Create(req models.CreatePekerjaanRequest) (*models.PekerjaanAlumni, error) {
	var newPekerjaan models.PekerjaanAlumni
	query := `
		INSERT INTO pekerjaan_alumni (alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, tanggal_mulai_kerja, gaji_range, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at
	`
	err := r.db.QueryRow(
		query,
		req.AlumniID, req.NamaPerusahaan, req.PosisiJabatan, req.BidangIndustri, req.LokasiKerja,
		req.TanggalMulaiKerja, req.GajiRange, req.TanggalSelesaiKerja, req.StatusPekerjaan,
		req.DeskripsiPekerjaan, time.Now(), time.Now(),
	).Scan(
		&newPekerjaan.ID, &newPekerjaan.AlumniID, &newPekerjaan.NamaPerusahaan, &newPekerjaan.PosisiJabatan,
		&newPekerjaan.BidangIndustri, &newPekerjaan.LokasiKerja, &newPekerjaan.GajiRange, &newPekerjaan.TanggalMulaiKerja,
		&newPekerjaan.TanggalSelesaiKerja, &newPekerjaan.StatusPekerjaan, &newPekerjaan.DeskripsiPekerjaan,
		&newPekerjaan.CreatedAt, &newPekerjaan.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &newPekerjaan, nil
}

func (r *pekerjaanRepository) Update(id int, req models.UpdatePekerjaanRequest) (*models.PekerjaanAlumni, error) {
	var updatedPekerjaan models.PekerjaanAlumni
	query := `
		UPDATE pekerjaan_alumni
		SET alumni_id = $1, nama_perusahaan = $2, posisi_jabatan = $3, bidang_industri = $4, lokasi_kerja = $5, tanggal_mulai_kerja = $6, gaji_range = $7, tanggal_selesai_kerja = $8, status_pekerjaan = $9, deskripsi_pekerjaan = $10, updated_at = $11
		WHERE id = $12
		RETURNING id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at
	`
	err := r.db.QueryRow(
		query,
		req.AlumniID, req.NamaPerusahaan, req.PosisiJabatan, req.BidangIndustri, req.LokasiKerja,
		req.TanggalMulaiKerja, req.GajiRange, req.TanggalSelesaiKerja, req.StatusPekerjaan,
		req.DeskripsiPekerjaan, time.Now(), id,
	).Scan(
		&updatedPekerjaan.ID, &updatedPekerjaan.AlumniID, &updatedPekerjaan.NamaPerusahaan, &updatedPekerjaan.PosisiJabatan,
		&updatedPekerjaan.BidangIndustri, &updatedPekerjaan.LokasiKerja, &updatedPekerjaan.GajiRange, &updatedPekerjaan.TanggalMulaiKerja,
		&updatedPekerjaan.TanggalSelesaiKerja, &updatedPekerjaan.StatusPekerjaan, &updatedPekerjaan.DeskripsiPekerjaan,
		&updatedPekerjaan.CreatedAt, &updatedPekerjaan.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &updatedPekerjaan, nil
}

func (r *pekerjaanRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM pekerjaan_alumni WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
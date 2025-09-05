package models

import (
	"database/sql"
	"time"
)

type PekerjaanAlumni struct {
	ID                  int64          `json:"id" db:"id"`
	AlumniID            int64          `json:"alumni_id" db:"alumni_id"`
	NamaPerusahaan      string         `json:"nama_perusahaan" db:"nama_perusahaan"`
	PosisiJabatan       string         `json:"posisi_jabatan" db:"posisi_jabatan"`
	BidangIndustri      string         `json:"bidang_industri" db:"bidang_industri"`
	LokasiKerja         string         `json:"lokasi_kerja" db:"lokasi_kerja"`
	GajiRange           sql.NullString `json:"gaji_range" db:"gaji_range"`                 
	TanggalMulaiKerja   time.Time      `json:"tanggal_mulai_kerja" db:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja sql.NullTime   `json:"tanggal_selesai_kerja" db:"tanggal_selesai_kerja"` 
	StatusPekerjaan     string         `json:"status_pekerjaan" db:"status_pekerjaan"`
	DeskripsiPekerjaan  sql.NullString `json:"deskripsi_pekerjaan" db:"deskripsi_pekerjaan"` 
	CreatedAt           time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at" db:"updated_at"`
}

type CreatePekerjaanRequest struct {
	AlumniID int64 `json:"alumni_id" validate:"required"`
	NamaPerusahaan    string `json:"nama_perusahaan" validate:"required,min=2"`
	PosisiJabatan     string `json:"posisi_jabatan" validate:"required,min=2"`
	BidangIndustri    string `json:"bidang_industri" validate:"required"`
	LokasiKerja       string `json:"lokasi_kerja" validate:"required"`
	TanggalMulaiKerja string `json:"tanggal_mulai_kerja" validate:"required,datetime=2006-01-02"`
	GajiRange           *string `json:"gaji_range,omitempty"`
	TanggalSelesaiKerja *string `json:"tanggal_selesai_kerja,omitempty" validate:"omitempty,datetime=2006-01-02"`
	StatusPekerjaan     *string `json:"status_pekerjaan,omitempty" validate:"omitempty,oneof=aktif selesai resigned"`
	DeskripsiPekerjaan  *string `json:"deskripsi_pekerjaan,omitempty"`
}

type UpdatePekerjaanRequest struct {
	AlumniID 			int64 `json:"alumni_id" validate:"required"`
	NamaPerusahaan    	string `json:"nama_perusahaan,omitempty" validate:"omitempty,min=2"`
	PosisiJabatan     	string `json:"posisi_jabatan,omitempty" validate:"omitempty,min=2"`
	BidangIndustri    	string `json:"bidang_industri,omitempty"`
	LokasiKerja       	string `json:"lokasi_kerja,omitempty"`
	TanggalMulaiKerja 	string `json:"tanggal_mulai_kerja,omitempty" validate:"omitempty,datetime=2006-01-02"`
	GajiRange           *string `json:"gaji_range,omitempty"`
	TanggalSelesaiKerja *string `json:"tanggal_selesai_kerja,omitempty" validate:"omitempty,datetime=2006-01-02"`
	StatusPekerjaan     *string `json:"status_pekerjaan,omitempty" validate:"omitempty,oneof=aktif selesai resigned"`
	DeskripsiPekerjaan  *string `json:"deskripsi_pekerjaan,omitempty"`
}
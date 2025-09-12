package models

import (
	"time"
) 

type Alumni struct{
	ID 	int `json:"id"`
	NIM string `json:"nim"`
	Nama string `json:"nama"`
	Jurusan string `json:"jurusan"`
	Angkatan int `json:"angkatan"`
	TahunLulus int `json:"tahun_lulus"`
	Email string `json:"email"`
	NoTelepon string `json:"no_telepon"`
	Alamat string `json:"alamat"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AlumniPekerjaan struct {
	ID                     int64  `json:"id"`
	Nama                   string `json:"nama"`
	Jurusan                string `json:"jurusan"`
	TahunLulus             int    `json:"tahun_lulus"`
	BidangIndustri         string `json:"bidang_industri"`
	NamaPerusahaan         string `json:"nama_perusahaan"`
	PosisiJabatan          string `json:"posisi_jabatan"`
	GajiRange              string `json:"gaji_range"`
	JumlahAlumniPerTahun   int    `json:"jumlah_alumni_per_tahun"`
}

type CreateAlumniRequest struct{
	NIM string `json:"nim"`
	Nama string `json:"nama"`
	Jurusan string `json:"jurusan"`
	Angkatan int `json:"angkatan"`
	TahunLulus int `json:"tahun_lulus"`
	Email string `json:"email"`
	NoTelepon string `json:"no_telepon"`
	Alamat string `json:"alamat"`
}

type UpdateAlumniRequest struct{
	NIM string `json:"nim"`
	Nama string `json:"nama"`
	Jurusan string `json:"jurusan"`
	Angkatan int `json:"angkatan"`
	TahunLulus int `json:"tahun_lulus"`
	Email string `json:"email"`
	NoTelepon string `json:"no_telepon"`
	Alamat string `json:"alamat"`
}
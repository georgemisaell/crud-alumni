package service

import (
	"crud-alumni/app/models"
	"crud-alumni/app/repository"
)

// PekerjaanService sekarang mendefinisikan interface untuk logika bisnis pekerjaan.
// Tidak ada lagi *fiber.Ctx.
type PekerjaanService interface {
	GetAllPekerjaan() ([]models.PekerjaanAlumni, error)
	GetPekerjaanByID(id int) (*models.PekerjaanAlumni, error)
	CreatePekerjaan(req models.CreatePekerjaanRequest) (*models.PekerjaanAlumni, error)
	UpdatePekerjaan(id int, req models.UpdatePekerjaanRequest) (*models.PekerjaanAlumni, error)
	DeletePekerjaan(id int) error
}

// pekerjaanService adalah implementasi konkretnya.
type pekerjaanService struct {
	repo repository.PekerjaanRepository
}

// NewPekerjaanService tidak berubah.
func NewPekerjaanService(repo repository.PekerjaanRepository) PekerjaanService {
	return &pekerjaanService{repo: repo}
}

// Method-method di bawah ini menjadi lebih fokus pada alur kerja data.

func (s *pekerjaanService) GetAllPekerjaan() ([]models.PekerjaanAlumni, error) {
	return s.repo.FindAll()
}

func (s *pekerjaanService) GetPekerjaanByID(id int) (*models.PekerjaanAlumni, error) {
	return s.repo.FindByID(id)
}

func (s *pekerjaanService) CreatePekerjaan(req models.CreatePekerjaanRequest) (*models.PekerjaanAlumni, error) {
	// Di sini tempatnya validasi bisnis yang lebih rumit jika diperlukan.
	return s.repo.Create(req)
}

func (s *pekerjaanService) UpdatePekerjaan(id int, req models.UpdatePekerjaanRequest) (*models.PekerjaanAlumni, error) {
	return s.repo.Update(id, req)
}

func (s *pekerjaanService) DeletePekerjaan(id int) error {
	return s.repo.Delete(id)
}
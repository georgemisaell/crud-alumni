package service

import (
	"crud-alumni/app/models"
	"crud-alumni/app/repository"
)

type PekerjaanService interface {
	GetAllPekerjaan() ([]models.PekerjaanAlumni, error)
	GetPekerjaanByID(id int) (*models.PekerjaanAlumni, error)
	CreatePekerjaan(req models.CreatePekerjaanRequest) (*models.PekerjaanAlumni, error)
	UpdatePekerjaan(id int, req models.UpdatePekerjaanRequest) (*models.PekerjaanAlumni, error)
	DeletePekerjaan(id int) error
}

type pekerjaanService struct {
	repo repository.PekerjaanRepository
}

func NewPekerjaanService(repo repository.PekerjaanRepository) PekerjaanService {
	return &pekerjaanService{repo: repo}
}

func (s *pekerjaanService) GetAllPekerjaan() ([]models.PekerjaanAlumni, error) {
	return s.repo.FindAll()
}

func (s *pekerjaanService) GetPekerjaanByID(id int) (*models.PekerjaanAlumni, error) {
	return s.repo.FindByID(id)
}

func (s *pekerjaanService) CreatePekerjaan(req models.CreatePekerjaanRequest) (*models.PekerjaanAlumni, error) {
	return s.repo.Create(req)
}

func (s *pekerjaanService) UpdatePekerjaan(id int, req models.UpdatePekerjaanRequest) (*models.PekerjaanAlumni, error) {
	return s.repo.Update(id, req)
}

func (s *pekerjaanService) DeletePekerjaan(id int) error {
	return s.repo.Delete(id)
}
package service

import (
	"crud-alumni/app/models"
	"crud-alumni/app/repository"
	"strings"
)

type PekerjaanService interface {
	GetAllPekerjaan() ([]models.PekerjaanAlumni, error)
	GetPekerjaanByID(id int) (*models.PekerjaanAlumni, error)
	GetWithPagination(search, sortBy, order string, page, limit int) ([]models.PekerjaanAlumni, int, error)
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

func (s *pekerjaanService) GetWithPagination(search, sortBy, order string, page, limit int) ([]models.PekerjaanAlumni, int, error) {
	sortByWhitelist := map[string]bool{"id": true, "nama_perusahaan": true, "bidang_industri": true, "created_at": true}
	if !sortByWhitelist[sortBy] {
		sortBy = "id"
	}

	if strings.ToLower(order) != "desc" {
		order = "asc"
	}

	offset := (page - 1) * limit

	pekerjaan, err := s.repo.FindWithPagination(search, sortBy, order, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count(search)
	if err != nil {
		return nil, 0, err
	}

	return pekerjaan, total, nil
}
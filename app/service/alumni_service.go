package service

import (
	"crud-alumni/app/models"
	"crud-alumni/app/repository"
	"strings"
)

type AlumniService interface {
	GetAllAlumni() ([]models.Alumni, error)
	GetAlumniByYear(year int) ([]models.AlumniPekerjaan, error)
	GetAlumniByID(id int) (*models.Alumni, error)
	GetAlumniWithPagination(search, sortBy, order string, page, limit int) ([]models.Alumni, int, error)
	CreateAlumni(req models.CreateAlumniRequest) (*models.Alumni, error)
	UpdateAlumni(id int, req models.UpdateAlumniRequest) (*models.Alumni, error)
	DeleteAlumni(id int) error
}

type alumniService struct {
	repo repository.AlumniRepository
}

func NewAlumniService(repo repository.AlumniRepository) AlumniService {
	return &alumniService{repo: repo}
}

func (s *alumniService) GetAllAlumni() ([]models.Alumni, error) {
	return s.repo.FindAll()
}

func (s *alumniService) GetAlumniByYear(year int) ([]models.AlumniPekerjaan, error) {
	return s.repo.FindByYear(year)
}

func (s *alumniService) GetAlumniByID(id int) (*models.Alumni, error) {
	return s.repo.FindByID(id)
}

func (s *alumniService) CreateAlumni(req models.CreateAlumniRequest) (*models.Alumni, error) {
	return s.repo.Create(req)
}

func (s *alumniService) UpdateAlumni(id int, req models.UpdateAlumniRequest) (*models.Alumni, error) {
	return s.repo.Update(id, req)
}

func (s *alumniService) DeleteAlumni(id int) error {
	return s.repo.Delete(id)
}

func (s *alumniService) GetAlumniWithPagination(search, sortBy, order string, page, limit int) ([]models.Alumni, int, error) {
	sortByWhitelist := map[string]bool{"id": true, "nim": true, "nama": true, "created_at": true}
	if !sortByWhitelist[sortBy] {
		sortBy = "id"
	}

	if strings.ToLower(order) != "desc" {
		order = "asc"
	}

	offset := (page - 1) * limit

	alumni, err := s.repo.FindWithPagination(search, sortBy, order, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count(search)
	if err != nil {
		return nil, 0, err
	}

	return alumni, total, nil
}
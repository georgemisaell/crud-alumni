package service

import (
	"crud-alumni/app/models"
	"crud-alumni/app/repository"
)

// AlumniService sekarang mendefinisikan interface untuk logika bisnis murni.
// Tidak ada lagi *fiber.Ctx di sini.
type AlumniService interface {
	GetAllAlumni() ([]models.Alumni, error)
	GetAlumniByYear(year int) ([]models.AlumniPekerjaan, error)
	GetAlumniByID(id int) (*models.Alumni, error)
	CreateAlumni(req models.CreateAlumniRequest) (*models.Alumni, error)
	UpdateAlumni(id int, req models.UpdateAlumniRequest) (*models.Alumni, error)
	DeleteAlumni(id int) error
}

// alumniService tetap menjadi implementasi konkret.
type alumniService struct {
	repo repository.AlumniRepository
}

// NewAlumniService tidak berubah, tetap menerima repository.
func NewAlumniService(repo repository.AlumniRepository) AlumniService {
	return &alumniService{repo: repo}
}

// Method-method di bawah ini sekarang lebih sederhana dan fokus pada logika.

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
	// Di masa depan, validasi bisnis yang kompleks bisa ditaruh di sini.
	// Contoh:
	// if s.repo.IsNIMExists(req.NIM) {
	//     return nil, errors.New("NIM sudah terdaftar")
	// }
	return s.repo.Create(req)
}

func (s *alumniService) UpdateAlumni(id int, req models.UpdateAlumniRequest) (*models.Alumni, error) {
	return s.repo.Update(id, req)
}

func (s *alumniService) DeleteAlumni(id int) error {
	return s.repo.Delete(id)
}
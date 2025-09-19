// app/service/auth_service.go
package service

import (
	"crud-alumni/app/models"
	"crud-alumni/app/repository"
	"crud-alumni/utils"
	"database/sql"
	"errors"
)

type AuthService interface {
	Login(req models.LoginRequest) (*models.LoginResponse, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	// 1. Cari user via repository
	user, passwordHash, err := s.userRepo.FindByUsernameOrEmail(req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("username atau password salah")
		}
		// Untuk error lainnya, kita anggap sebagai internal error
		return nil, errors.New("terjadi kesalahan pada server")
	}

	// 2. Cek password
	if !utils.CheckPassword(req.Password, passwordHash) {
		return nil, errors.New("username atau password salah")
	}

	// 3. Generate token jika password cocok
	token, err := utils.GenerateToken(*user)
	if err != nil {
		return nil, errors.New("gagal membuat token")
	}

	// 4. Siapkan response
	response := &models.LoginResponse{
		User:  *user,
		Token: token,
	}

	return response, nil
}
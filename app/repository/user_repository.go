// app/repository/user_repository.go
package repository

import (
	"crud-alumni/app/models"
	"database/sql"
)

// Definisikan interface untuk user repository agar mudah untuk di-mock saat testing
type UserRepository interface {
	FindByUsernameOrEmail(usernameOrEmail string) (*models.User, string, error)
}

// Struct implementasi dari interface
type userRepository struct {
	db *sql.DB
}

// Constructor untuk membuat instance baru dari userRepository
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// Method untuk mencari user berdasarkan username atau email
func (r *userRepository) FindByUsernameOrEmail(usernameOrEmail string) (*models.User, string, error) {
	var user models.User
	var passwordHash string

	query := `
        SELECT id, username, email, password_hash, role, created_at
        FROM users
        WHERE username = $1 OR email = $1
    `
	err := r.db.QueryRow(query, usernameOrEmail).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&passwordHash,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, "", err // Biarkan service yang menangani jenis errornya
	}

	return &user, passwordHash, nil
}
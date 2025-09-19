// app/repository/user_repository.go
package repository

import (
	"crud-alumni/app/models"
	"database/sql"
)

type UserRepository interface {
	FindByUsernameOrEmail(usernameOrEmail string) (*models.User, string, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

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
		return nil, "", err
	}

	return &user, passwordHash, nil
}
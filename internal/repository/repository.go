package repository

import (
	"database/sql"
	"fmt"
	"login-with-oauth/internal/logger"
	"login-with-oauth/internal/models"
)

// UserRepository is the interface for the user repository
type UserRepository interface {
	CreateUser(user models.User) (*models.User, error)
	GetUserByID(id string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
}

// UserRepositoryImpl is the implementation of the UserRepository interface
type UserRepositoryImpl struct {
	db *sql.DB
}

// NewUserRepository creates a new instance of the UserRepository
func NewUserRepository(db *sql.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

// CreateUser creates a new user
func (r *UserRepositoryImpl) CreateUser(user models.User) (*models.User, error) {
	fmt.Printf("Repository: Creating user with data: %+v\n", user)

	// PostgreSQL upsert syntax using ON CONFLICT
	query := `
		INSERT INTO users (id, username, email, avatar_url, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (email) DO UPDATE SET 
			email = EXCLUDED.email,
			updated_at = EXCLUDED.updated_at
		RETURNING *`

	// For PostgreSQL, use QueryRow to get the returned row
	var savedUser models.User
	err := r.db.QueryRow(query,
		user.ID,
		user.Username,
		user.Email,
		user.AvatarURL,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(
		&savedUser.ID,
		&savedUser.Username,
		&savedUser.Email,
		&savedUser.AvatarURL,
		&savedUser.CreatedAt,
		&savedUser.UpdatedAt,
	)

	if err != nil {
		logger.Log.Error("Failed to execute insert query: " + err.Error())
		return nil, fmt.Errorf("failed to execute insert query: %v", err)
	}

	return &savedUser, nil
}

// GetUserByID retrieves a user by their ID
func (r *UserRepositoryImpl) GetUserByID(id string) (*models.User, error) {
	// TODO: Implement user retrieval logic
	query := "SELECT id, username, email, avatar_url, created_at, updated_at FROM users WHERE id = ?"
	row := r.db.QueryRow(query, id)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail retrieves a user by their email
func (r *UserRepositoryImpl) GetUserByEmail(email string) (*models.User, error) {
	// TODO: Implement user retrieval logic
	query := "SELECT id, username, email, avatar_url, created_at, updated_at FROM users WHERE email = ?"
	row := r.db.QueryRow(query, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

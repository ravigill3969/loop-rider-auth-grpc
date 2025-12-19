package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"ravigill/rider-grpc-server/internal/models"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	query := `INSERT INTO users (id, email, full_name, password, phone_number, birth_month, birth_year, updated_at, created_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	user.CreatedAt = time.Now().Unix()
	user.UpdatedAt = time.Now().Unix()

	err := r.db.QueryRowContext(ctx, query, user.ID, user.Email, user.FullName, user.Password, user.PhoneNumber, user.BirthMonth, user.BirthYear, user.UpdatedAt, user.CreatedAt).Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, email, full_name, password, phone_number, birth_month, birth_year, updated_at, created_at FROM users WHERE email = $1`
	row := r.db.QueryRowContext(ctx, query, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.FullName, &user.Password, &user.PhoneNumber, &user.BirthMonth, &user.BirthYear, &user.UpdatedAt, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	query := `SELECT id, email, full_name, password, phone_number, birth_month, birth_year, updated_at, created_at FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.FullName, &user.Password, &user.PhoneNumber, &user.BirthMonth, &user.BirthYear, &user.UpdatedAt, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

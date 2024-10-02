package userRepository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/codepnw/ticket-api/models"
	"github.com/jmoiron/sqlx"
)

type IUserRepository interface {
	InsertUser(ctx context.Context, req *models.UserRegisterReq, isAdmin bool) (*models.UserPassport, error)
	GetProfile(userId uint) (*models.User, error)
	FindOneUserByEmail(email string) (*models.User, error)
}

type userRepository struct {
	db *sqlx.DB
	id string
}

func NewUserRepository(db *sqlx.DB) IUserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) createUser(ctx context.Context, req *models.UserRegisterReq) error {
	query := `
		INSERT INTO users (email, password, username, role)
		VALUES ($1, $2, $3, 1)
		RETURNING id;
	`
	err := r.db.QueryRowContext(ctx, query, req.Email, req.Password, req.Username).Scan(&r.id)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) createAdmin(ctx context.Context, req *models.UserRegisterReq) error {
	query := `
		INSERT INTO users (email, password, username, role)
		VALUES ($1, $2, $3, 2)
		RETURNING id;
	`
	err := r.db.QueryRowContext(ctx, query, req.Email, req.Password, req.Username).Scan(&r.id)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) InsertUser(ctx context.Context, req *models.UserRegisterReq, isAdmin bool) (*models.UserPassport, error) {
	if isAdmin {
		if err := r.createAdmin(ctx, req); err != nil {
			return nil, err
		}
	} else {
		if err := r.createUser(ctx, req); err != nil {
			return nil, err
		}
	}

	query := `
		SELECT json_build_object('user', t, 'token', NULL)
		FROM (
			SELECT u.id, u.email, u.username, u.role, u.created_at, u.updated_at
			FROM users u
			WHERE u.id = $1
		) AS t
	`
	data := make([]byte, 0)
	if err := r.db.Get(&data, query, r.id); err != nil {
		return nil, fmt.Errorf("get user failed: %w", err)
	}

	user := new(models.UserPassport)
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("unmarshal user failed: %w", err)
	}
	return user, nil
}

func (r *userRepository) GetProfile(userId uint) (*models.User, error) {
	query := `
		SELECT id, email, username, role
		FROM users
		WHERE id = $1;
	`
	profile := new(models.User)

	if err := r.db.Get(profile, query, userId); err != nil {
		return nil, err
	}
	return profile, nil
}

func (r *userRepository) FindOneUserByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, password, username, role
		FROM users
		WHERE email = $1;
	`
	user := new(models.User)

	if err := r.db.Get(user, query, email); err != nil {
		return nil, err
	}
	return user, nil
}

package userRepository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/codepnw/ticket-api/models"
	"github.com/jmoiron/sqlx"
)

type IUserRepository interface {
	InsertUser(ctx context.Context, req *models.UserRegisterReq, isAdmin bool) (*models.UserPassport, error)
	GetProfile(userId uint) (*models.User, error)
	FindOneUserByEmail(email string) (*models.UserCredentialCheck, error)
	InsertOauth(req *models.UserPassport) error
	FindOneOauth(refreshToken string) (*models.Oauth, error)
	UpdateOauth(req *models.UserToken) error
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

func (r *userRepository) FindOneUserByEmail(email string) (*models.UserCredentialCheck, error) {
	query := `
		SELECT id, email, password, username, role
		FROM users
		WHERE email = $1;
	`
	user := new(models.UserCredentialCheck)

	if err := r.db.Get(user, query, email); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) InsertOauth(req *models.UserPassport) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `
		INSERT INTO "oauth" (
			"user_id",
			"refresh_token",
			"access_token"
		)
		VALUES ($1, $2, $3)
		RETURNING "id";`

	err := r.db.QueryRowContext(
		ctx,
		query,
		req.User.ID,
		req.Token.RefreshToken,
		req.Token.AccessToken,
	).Scan(&req.Token.Id)

	if err != nil {
		return fmt.Errorf("insert oauth failed: %v", err)
	}
	return nil
}

func (r *userRepository) FindOneOauth(refreshToken string) (*models.Oauth, error) {
	query := `
		SELECT id, user_id
		FROM oauth
		WHERE refresh_token = $1;
	`
	oauth := new(models.Oauth)

	if err := r.db.Get(oauth, query, refreshToken); err != nil {
		return nil, err
	}
	return oauth, nil
}

func (r *userRepository) UpdateOauth(req *models.UserToken) error {
	query := `
		UPDATE oauth SET
			access_token = :access_token,
			refresh_token = :refresh_token
		WHERE id = :id;
	`
	if _, err := r.db.NamedExecContext(context.Background(), query, req); err != nil {
		return err
	}
	return nil
}

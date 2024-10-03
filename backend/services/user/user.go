package userService

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/codepnw/ticket-api/cmd/config"
	"github.com/codepnw/ticket-api/models"
	"github.com/codepnw/ticket-api/pkg/auth"
	userRepository "github.com/codepnw/ticket-api/repositories/user"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	CreateAdmin(req *models.UserRegisterReq) (*models.UserPassport, error)
	CreateUser(req *models.UserRegisterReq) (*models.UserPassport, error)
	GetProfile(userId uint) (*models.User, error)
	GetPassport(req *models.UserCredential) (*models.UserPassport, error)
	RefreshPassport(req *models.UserRefreshCredential) (*models.UserPassport, error)
}

type userService struct {
	cfg      config.EnvConfig
	userRepo userRepository.IUserRepository
}

func NewUserService(cfg config.EnvConfig, userRepo userRepository.IUserRepository) IUserService {
	return &userService{cfg: cfg, userRepo: userRepo}
}

func (s *userService) CreateAdmin(req *models.UserRegisterReq) (*models.UserPassport, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if !req.IsEmail() {
		return nil, fmt.Errorf("invalid email address")
	}

	if err := req.BcryptHashing(); err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed hashing password")
	}

	result, err := s.userRepo.InsertUser(ctx, req, true)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("create admin failed")
	}

	return result, nil
}

func (s *userService) CreateUser(req *models.UserRegisterReq) (*models.UserPassport, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if !req.IsEmail() {
		return nil, fmt.Errorf("invalid email address")
	}

	if err := req.BcryptHashing(); err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed hashing password")
	}

	result, err := s.userRepo.InsertUser(ctx, req, false)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("create user failed")
	}

	return result, nil
}

func (s *userService) GetProfile(userId uint) (*models.User, error) {
	user, err := s.userRepo.GetProfile(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user_id not found")
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) GetPassport(req *models.UserCredential) (*models.UserPassport, error) {
	user, err := s.userRepo.FindOneUserByEmail(req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("email or password is invalid")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("email or password is invalid")
	}

	accessToken, _ := auth.NewTicketAuth(auth.Access, s.cfg, &models.UserClaims{
		Id:     user.ID,
		RoleId: user.Role,
	})

	refreshToken, _ := auth.NewTicketAuth(auth.Refresh, s.cfg, &models.UserClaims{
		Id:     user.ID,
		RoleId: user.Role,
	})

	ssAccessToken, err := accessToken.SignToken()
	if err != nil {
		return nil, err
	}

	ssRefreshToken, err := refreshToken.SignToken()
	if err != nil {
		return nil, err
	}

	passport := &models.UserPassport{
		User: &models.User{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
			Role:     user.Role,
		},
		Token: &models.UserToken{
			AccessToken:  ssAccessToken,
			RefreshToken: ssRefreshToken,
		},
	}

	if err := s.userRepo.InsertOauth(passport); err != nil {
		return nil, err
	}
	return passport, nil
}

func (s *userService) RefreshPassport(req *models.UserRefreshCredential) (*models.UserPassport, error) {
	claims, err := auth.ParseToken(s.cfg, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	oauth, err := s.userRepo.FindOneOauth(req.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("refresh_token not found")
	}

	profile, err := s.userRepo.GetProfile(oauth.UserId)
	if err != nil {
		return nil, fmt.Errorf("profile not found")
	}

	newClaims := &models.UserClaims{
		Id:     profile.ID,
		RoleId: profile.Role,
	}

	accessToken, err := auth.NewTicketAuth(auth.Access, s.cfg, newClaims)
	if err != nil {
		return nil, err
	}

	refreshToken := auth.RepeatToken(s.cfg, newClaims, claims.ExpiresAt.Unix())
	ssAccessToken, err := accessToken.SignToken()
	if err != nil {
		return nil, err
	}

	passport := &models.UserPassport{
		User: profile,
		Token: &models.UserToken{
			Id:           oauth.Id,
			AccessToken:  ssAccessToken,
			RefreshToken: refreshToken,
		},
	}

	if err := s.userRepo.UpdateOauth(passport.Token); err != nil {
		return nil, fmt.Errorf("update oauth failed: %w", err)
	}
	return passport, nil
}

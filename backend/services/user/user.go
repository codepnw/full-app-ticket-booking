package userService

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/codepnw/ticket-api/models"
	userRepository "github.com/codepnw/ticket-api/repositories/user"
)

type IUserService interface {
	CreateAdmin(req *models.UserRegisterReq) (*models.UserPassport, error)
	CreateUser(req *models.UserRegisterReq) (*models.UserPassport, error)
	GetProfile(userId uint) (*models.User, error)
}

type userService struct {
	userRepo userRepository.IUserRepository
}

func NewUserService(userRepo userRepository.IUserRepository) IUserService {
	return &userService{userRepo: userRepo}
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

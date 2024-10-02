package models

import (
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint   `db:"id" json:"id"`
	Email     string `db:"email" json:"email"`
	Username  string `db:"username" json:"username"`
	Role      int    `db:"role" json:"role"`
	CreatedAt string `db:"created_at" json:"createdAt"`
	UpdatedAt string `db:"updated_at" json:"updatedAt"`
}

type UserRegisterReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

func (req *UserRegisterReq) BcryptHashing() error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return err
	}
	req.Password = string(hashed)
	return nil
}

func (req *UserRegisterReq) IsEmail() bool {
	match, err := regexp.MatchString(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, req.Email)
	if err != nil {
		return false
	}
	return match
}

type UserCredential struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserPassport struct {
	User  *User      `json:"user"`
	Token *UserToken `json:"token"`
}

type UserToken struct {
	Id           string `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

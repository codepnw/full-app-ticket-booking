package auth

import (
	"fmt"
	"math"
	"time"

	"github.com/codepnw/ticket-api/cmd/config"
	"github.com/codepnw/ticket-api/models"
	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	Access  TokenType = "access"
	Refresh TokenType = "refresh"
	Admin   TokenType = "admin"
	ApiKey  TokenType = "apiKey"
)

const issuerStr string = "ticket-booking-api"

func NewTicketAuth(tokenType TokenType, cfg config.EnvConfig, claims *models.UserClaims) (ITicketAuth, error) {
	switch tokenType {
	case Access:
		return newAccessToken(cfg, claims), nil
	case Refresh:
		return newRefreshToken(cfg, claims), nil
	case Admin:
		return newAdminToken(cfg), nil
	case ApiKey:
		return newApiKey(cfg), nil
	default:
		return nil, fmt.Errorf("unknow token type")
	}
}

type ITicketAuth interface {
	SignToken() (string, error)
}

type ITicketAdmin interface {
	SignToken() (string, error)
}

type ITicketApiKey interface {
	SignToken() (string, error)
}

type ticketAuth struct {
	cfg       config.EnvConfig
	mapClaims *ticketMapClaims
}

type ticketAdminKey struct {
	*ticketAuth
}

type ticketApiKey struct {
	*ticketAuth
}

type ticketMapClaims struct {
	Claims *models.UserClaims `json:"claims"`
	jwt.RegisteredClaims
}

func (a *ticketAuth) SignToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, err := token.SignedString([]byte(a.cfg.JWTSecretKey))
	if err != nil {
		return "", fmt.Errorf("sign secret token failed: %v", err)
	}
	return ss, nil
}

func (a *ticketAdminKey) SignToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, err := token.SignedString([]byte(a.cfg.JWTAdminKey))
	if err != nil {
		return "", fmt.Errorf("sign admin token failed: %v", err)
	}
	return ss, nil
}

func (a *ticketApiKey) SignToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, err := token.SignedString([]byte(a.cfg.JWTApiKey))
	if err != nil {
		return "", fmt.Errorf("sign api_key failed: %v", err)
	}
	return ss, nil
}

func newAccessToken(cfg config.EnvConfig, cliams *models.UserClaims) ITicketAuth {
	return &ticketAuth{
		cfg: cfg,
		mapClaims: &ticketMapClaims{
			Claims: cliams,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    issuerStr,
				Subject:   "access-token",
				Audience:  []string{"admin", "customer"},
				ExpiresAt: jwtTimeDurationCal(cfg.JWTExp),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func newRefreshToken(cfg config.EnvConfig, cliams *models.UserClaims) ITicketAuth {
	return &ticketAuth{
		cfg: cfg,
		mapClaims: &ticketMapClaims{
			Claims: cliams,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    issuerStr,
				Subject:   "refresh-token",
				Audience:  []string{"admin", "customer"},
				ExpiresAt: jwtTimeDurationCal(cfg.JWTExp),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func newAdminToken(cfg config.EnvConfig) ITicketAuth {
	return &ticketAdminKey{
		ticketAuth: &ticketAuth{
			cfg: cfg,
			mapClaims: &ticketMapClaims{
				Claims: nil,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    issuerStr,
					Subject:   "admin-token",
					Audience:  []string{"admin"},
					ExpiresAt: jwtTimeDurationCal(300), // 5 min
					NotBefore: jwt.NewNumericDate(time.Now()),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		},
	}
}

func newApiKey(cfg config.EnvConfig) ITicketAuth {
	return &ticketApiKey{
		ticketAuth: &ticketAuth{
			cfg: cfg,
			mapClaims: &ticketMapClaims{
				Claims: nil,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    issuerStr,
					Subject:   "api-key",
					Audience:  []string{"admin", "customer"},
					ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(1, 0, 0)), // 1 year
					NotBefore: jwt.NewNumericDate(time.Now()),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		},
	}
}

func jwtTimeDurationCal(t int) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Duration(int64(t) * int64(math.Pow10(9)))))
}

func jwtTimeRepeatAdapter(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(t, 0))
}

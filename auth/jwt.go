package auth

import (
	"time"
	"urfunavigator/index/models"

	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	Secret    []byte
	ExpiresIn time.Duration
}

type Claims struct {
	jwt.RegisteredClaims
	UserID string          `json:"userId"`
	Login  string          `json:"login"`
	Role   models.UserRole `json:"role"`
}

func GenerateToken(user *models.User, cfg Config) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(cfg.ExpiresIn)

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   user.Id.Hex(),
		},
		UserID: user.Id.Hex(),
		Login:  user.Login,
		Role:   user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(cfg.Secret)
	if err != nil {
		return "", time.Time{}, err
	}

	return signed, expiresAt, nil
}

package service

import (
	"errors"
	"github.com/MCPutro/go-management-project/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWTService interface {
	GenerateToken(userID int64, email string) (string, error)
	ValidateToken(token string) (int64, error)
}
type jwtService struct {
	secretKey []byte
	expiresIn time.Time
}

func NewJwtService(config *config.JwtConfig) JWTService {
	return &jwtService{
		secretKey: []byte(config.SecretKey),
		expiresIn: time.Now().Add(time.Duration(config.ExpirationInSecond) * time.Second),
	}
}

type claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func (j *jwtService) GenerateToken(userID int64, email string) (string, error) {
	claim := &claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(j.expiresIn),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *jwtService) ValidateToken(tokenString string) (int64, error) {
	claims := &claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Pastikan signing method sesuai
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.secretKey, nil
	})

	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	return claims.UserID, nil
}

package auth

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type TokenService interface {
	IssueJWT(userID, userName, realm string) (string, error)
}

type jwtTokenService struct {
	secret string
}

func NewJwtTokenService(secret string) TokenService {
	return &jwtTokenService{secret: secret}
}

func (service *jwtTokenService) IssueJWT(userID, userName, realm string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":       userID,
		"user_name": userName,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
		"realm":     realm,
	})

	tokenString, err := token.SignedString([]byte(service.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

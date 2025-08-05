package infrastructure

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)

type JWTService struct{

}
func NewJWTService()*JWTService{
	return &JWTService{}
}
var jwtSecret = []byte("access-secret")
var refreshSecret = []byte("refresh-secret")

func (j *JWTService)GenerateTokens(userID string) (string, string, error) {
	// Access Token
	atClaims := jwt.MapClaims{}
	atClaims["user_id"] = userID
	atClaims["exp"] = time.Now().Add(15 * time.Minute).Unix()
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	atString, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	// Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["user_id"] = userID
	rtClaims["exp"] = time.Now().Add(7 * 24 * time.Hour).Unix()
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	rtString, err := refreshToken.SignedString(refreshSecret)
	if err != nil {
		return "", "", err
	}

	return atString, rtString, nil
}

func (j *JWTService) ValidateAccessToken(tokenStr string) (string, error) {
	return j.ValidateToken(tokenStr, jwtSecret)
}

func (j *JWTService) ValidateRefreshToken(tokenStr string) (string, error) {
	return j.ValidateToken(tokenStr, refreshSecret)
}

func (j *JWTService) ValidateToken(tokenStr string, secret []byte) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["user_id"] == nil {
		return "", errors.New("invalid claims")
	}

	return claims["user_id"].(string), nil
}

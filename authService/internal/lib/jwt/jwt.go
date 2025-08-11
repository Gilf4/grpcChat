package jwt

import (
	"time"

	"github.com/Gilf4/grpcChat/auth/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
)

func NewToken(user *models.User, duration time.Duration, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	return token.SignedString(secret)
}

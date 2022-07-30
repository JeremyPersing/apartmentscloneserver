package utils

import (
	"os"
	"time"

	"github.com/kataras/iris/v12/middleware/jwt"
)

func CreateForgotPasswordToken(id uint, email string) (string, error) {
	signer := jwt.NewSigner(jwt.HS256, os.Getenv("EMAIL_TOKEN_SECRET"), 10*time.Minute)

	claims := ForgotPasswordToken{
		ID:    id,
		Email: email,
	}

	token, err := signer.Sign(claims)
	if err != nil {
		return "", err
	}

	return string(token), nil
}

type ForgotPasswordToken struct {
	ID    uint   `json:"ID"`
	Email string `json:"email"`
}

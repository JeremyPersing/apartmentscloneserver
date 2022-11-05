package utils

import (
	"apartments-clone-server/storage"
	"context"
	"os"
	"strconv"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
)

var bgContext = context.Background()

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

func CreateTokenPair(id uint) (*jwt.TokenPair, error) {
	accessTokenSigner := jwt.NewSigner(jwt.HS256, os.Getenv("ACCESS_TOKEN_SECRET"), 24*time.Hour)
	refreshTokenSigner := jwt.NewSigner(jwt.HS256, os.Getenv("REFRESH_TOKEN_SECRET"), 365*24*time.Hour)

	userID := strconv.FormatUint(uint64(id), 10)

	refreshClaims := jwt.Claims{Subject: userID}

	accessTokenClaims := AccessToken{
		ID: id,
	}

	accessToken, err := accessTokenSigner.Sign(accessTokenClaims)
	if err != nil {
		return nil, err
	}

	refreshToken, err := refreshTokenSigner.Sign(refreshClaims)
	if err != nil {
		return nil, err
	}

	var tokenPair jwt.TokenPair
	tokenPair.AccessToken = accessToken
	tokenPair.RefreshToken = refreshToken

	storage.Redis.Set(bgContext, string(refreshToken), "true", 365*24*time.Hour+5*time.Minute)

	return &tokenPair, nil
}

func RefreshToken(ctx iris.Context) {
	token := jwt.GetVerifiedToken(ctx)
	tokenStr := string(token.Token)
	validToken, tokenErr := storage.Redis.Get(bgContext, tokenStr).Result()

	if tokenErr != nil {
		CreateNotFound(ctx)
		return
	}

	if validToken != "true" {
		ctx.StatusCode(iris.StatusForbidden)
		return
	}

	storage.Redis.Del(bgContext, tokenStr)
	userID, parseErr := strconv.ParseUint(token.StandardClaims.Subject, 10, 32)
	if parseErr != nil {
		CreateInternalServerError(ctx)
		return
	}

	tokenPair, tokenPairErr := CreateTokenPair(uint(userID))
	if tokenPairErr != nil {
		CreateInternalServerError(ctx)
		return
	}

	ctx.JSON(iris.Map{
		"accessToken":  string(tokenPair.AccessToken),
		"refreshToken": string(tokenPair.RefreshToken),
	})
}

type ForgotPasswordToken struct {
	ID    uint   `json:"ID"`
	Email string `json:"email"`
}

type AccessToken struct {
	ID uint `json:"ID"`
}

type RefreshTokenInput struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

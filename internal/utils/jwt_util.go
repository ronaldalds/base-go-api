package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ronaldalds/base-go-api/internal/models"
	"github.com/ronaldalds/base-go-api/internal/settings"
)

type PayloadJwt struct {
	Token  string
	Claims JwtClaims
}
type JwtClaims struct {
	Sub uint `json:"sub"`
	Exp int  `json:"exp"`
	jwt.RegisteredClaims
}

func GenerateToken(user *models.User, expire time.Duration) (string, error) {
	location, err := time.LoadLocation(settings.Env.TimeZone)
	if err != nil {
		return "", fmt.Errorf("invalid timezone: %s", err.Error())
	}
	currentTime := time.Now().In(location)

	accessTokenExpirationTime := currentTime.Add(expire)

	accessClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"iat": currentTime.Unix(),
		"iss": "ponche",
		"exp": accessTokenExpirationTime.Unix(),
	})

	accessToken, err := accessClaims.SignedString([]byte(settings.Env.JwtSecret))
	if err != nil {
		return "", fmt.Errorf("could not sign access token string %v", err.Error())
	}

	return accessToken, nil
}

func GetJwtHeaderPayload(ctx *fiber.Ctx) (*PayloadJwt, error) {
	authHeader := ctx.Get("Authorization")
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		tokenSecret := settings.Env.JwtSecret
		return []byte(tokenSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid jwt token")
	}

	tokenDone := token.Claims.(*JwtClaims)
	jwt := &PayloadJwt{
		Token:  tokenString,
		Claims: *tokenDone,
	}

	return jwt, nil
}

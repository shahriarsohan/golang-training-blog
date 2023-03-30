package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateToken(ttl time.Duration, payload interface{}, sectetJWTKey string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	now := time.Now().UTC()
	claims := token.Claims.(jwt.MapClaims)

	claims["sub"] = payload
	claims["exp"] = now.Add(ttl).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	tokenString, err := token.SignedString([]byte(sectetJWTKey))

	if err != nil {
		return "", fmt.Errorf("generating JWT token failed %w", err)
	}

	return tokenString, nil
}

func ValidateToken(token string, signedJWTKey string) (interface{}, error) {
	tok, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); ok {
			return nil, fmt.Errorf("undexpected method %s", t.Header["alg"])
		}
		return []byte(signedJWTKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalidate token")
	}

	claims, ok := tok.Claims.(jwt.MapClaims)

	if !ok || !tok.Valid {
		return nil, fmt.Errorf("invalid token claim")
	}

	return claims["sub"], nil
}

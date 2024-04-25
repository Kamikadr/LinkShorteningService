package api

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Auth struct {
	accessTtl  time.Duration
	refreshTtl time.Duration
	signedKey  []byte
}

func NewAuth(accessTtl, refreshTtl time.Duration, signedKey string) *Auth {
	return &Auth{accessTtl, refreshTtl, []byte(signedKey)}
}

func (auth *Auth) GenerateToken(username string) (string, error) {
	return auth.generateToken(username, auth.accessTtl)
}

func (auth *Auth) GenerateRefreshToken(username string) (string, error) {
	return auth.generateToken(username, auth.refreshTtl)
}

func (auth *Auth) ParseToken(tokenString string) (string, error) {
	var claims jwt.RegisteredClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return auth.signedKey, nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}
	return claims.Subject, nil
}

func (auth *Auth) generateToken(username string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   username,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl))})
	return token.SignedString(auth.signedKey)
}

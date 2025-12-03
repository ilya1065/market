package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type JWT struct {
	secret        []byte
	JWTRefreshTTL int
	JWTAccessTTL  int
}

func NewJWT(secret string, JWTRefreshTTL, JWTAccessTTL int) *JWT {
	return &JWT{secret: []byte(secret),
		JWTAccessTTL:  JWTAccessTTL,
		JWTRefreshTTL: JWTRefreshTTL}
}

func (j *JWT) AccessGenerate(userID uint) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", userID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.JWTAccessTTL) * time.Second)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(j.secret)
}

func (j *JWT) RefreshGenerate(userID uint) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.JWTRefreshTTL) * time.Second)),
		Subject:   fmt.Sprintf("%d", userID),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        uuid.New().String(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(j.secret)
}

// TODO
func (j *JWT) Parse(tokenStr string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		return j.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if c, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return c, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}

package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type JWTValidator struct {
	secret []byte
}

func NewJWTValidator(secret string) *JWTValidator {
	return &JWTValidator{secret: []byte(secret)}
}

func (v *JWTValidator) Validate(_ context.Context, tokenOrBearer string) error {
	token := tokenOrBearer
	if strings.HasPrefix(strings.ToLower(tokenOrBearer), "bearer ") {
		token = strings.TrimSpace(tokenOrBearer[len("Bearer "):])
	}
	if token == "" {
		return errors.New("empty token")
	}

	_, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		// безопасно: ожидаем HS256
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return v.secret, nil
	})
	return err
}

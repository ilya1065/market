package http

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (h *Handler) AuthMiddleware(c *gin.Context) {
	fmt.Println("AuthMiddleware")
	authz := c.GetHeader("Authorization")
	fmt.Println(authz)
	if !strings.HasPrefix(strings.ToLower(authz), "bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
		return
	}
	tokenStr := strings.TrimSpace(authz[len("Bearer "):])

	// у UserService внутри есть h.users.jwt, но мы не хотим лезть в сервис
	claims, err := h.users.JWTParse(tokenStr) // добавим прокси-метод вниз
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}
	// subject хранит userID строкой
	uid, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil || !tokenValid(claims) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		return
	}
	c.Set("userID", uint(uid))
	c.Next()
}

func tokenValid(c *jwt.RegisteredClaims) bool {
	return c != nil &&
		c.ExpiresAt != nil &&
		c.ExpiresAt.Time.After(time.Now())
}

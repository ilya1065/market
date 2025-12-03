package http

import (
	"errors"
	"fmt"
	"marketplace/internal/auth"
	"marketplace/internal/service"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	users *service.UserService
}

func NewHandler(user *service.UserService) *Handler {
	return &Handler{users: user}
}

type registerReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ввод"})
	}
	users, token, err := h.users.Register(req.Email, req.Password, auth.HashPassword)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrEmailTaken) {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.SetCookie("Access_Token", token.Access, h.users.Jwt.JWTAccessTTL, "", os.Getenv("HOST"), true, true)
	c.SetCookie("Refresh_token", token.Refresh, h.users.Jwt.JWTRefreshTTL, "", os.Getenv("HOST"), true, true)
	c.JSON(http.StatusCreated, gin.H{"id": users.ID, "email": users.Email})
}

func (h *Handler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ввод"})
		return
	}
	token, err := h.users.Login(req.Email, req.Password, auth.CheckPassword)
	if err != nil {
		status := http.StatusUnauthorized
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.SetCookie("Access_Token", token.Access, h.users.Jwt.JWTAccessTTL, "", os.Getenv("HOST"), true, true)
	c.SetCookie("Refresh_token", token.Refresh, h.users.Jwt.JWTRefreshTTL, "", os.Getenv("HOST"), true, true)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handler) Me(c *gin.Context) {
	fmt.Println("me")
	id := c.GetUint("userID")
	users, err := h.users.GetByID(id)
	if err != nil || users == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": users.ID, "email": users.Email})
}

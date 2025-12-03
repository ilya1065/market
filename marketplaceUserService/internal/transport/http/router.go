package http

import "github.com/gin-gonic/gin"

func NewRouter(h *Handler) *gin.Engine {
	r := gin.Default()
	r.POST("/register", h.Register)
	//Get
	r.POST("/login", h.Login)

	auth := r.Group("/")
	auth.Use(h.AuthMiddleware)
	auth.GET("/me", h.Me)
	return r
}

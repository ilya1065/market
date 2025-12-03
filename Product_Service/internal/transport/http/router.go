package http

import "github.com/gin-gonic/gin"

func NewRouter(h *ProductHendler) *gin.Engine {
	r := gin.Default()
	r.POST("/create", h.Create)
	r.GET("/list", h.ListProducts)
	r.PUT("/update/:id", h.UpdateProduct)
	r.DELETE("/delete/:id", h.DeleteProduct)
	r.GET("/product", h.GetProductByID)
	return r
}

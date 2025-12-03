package http

import (
	"Product_Service/internal/entity"
	"Product_Service/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ProductHendler struct {
	product service.ProductService
}

func NewHendler(svc service.ProductService) *ProductHendler {
	return &ProductHendler{product: svc}
}

func (h ProductHendler) Create(c *gin.Context) {
	var product entity.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"masseg": "неверный ввод",
			"error":  err,
		})

	}
	if err := h.product.Create(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при создании продукта"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"product": product})
}

func (h ProductHendler) GetProductByID(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	idUint64, err := strconv.ParseUint(idStr, 10, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	id := uint(idUint64)
	product, err := h.product.GetProductByID(id)
	c.JSON(http.StatusOK, gin.H{"product": product})

}
func (h ProductHendler) ListProducts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	products, err := h.product.ListProducts(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"products": products})
}

func (h ProductHendler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	idInt, err := strconv.Atoi(idStr)
	id := uint(idInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	var product entity.Product
	if err = c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	res, err := h.product.UpdateProduct(id, &product)
	c.JSON(http.StatusOK, gin.H{"product": res})

}
func (h ProductHendler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	id := uint(idInt)
	err = h.product.DeleteProduct(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

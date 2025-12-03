package product

import (
	"Product_Service/internal/entity"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type productRepo struct {
	db *sqlx.DB
}

func (p productRepo) GetProductByID(id uint) (*entity.Product, error) {
	var product entity.Product
	err := p.db.Get(&product, "SELECT * FROM product WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить продукт из DB %w", err)
	}
	return &product, nil
}
func (p productRepo) Create(product *entity.Product) error {
	_, err := p.db.NamedExec(" INSERT INTO product (price, name, description, stock, image_url, category_id, created_at, updated_at) VALUES (:price, :name, :description, :stock, :image_url, :category_id, :created_at, :updated_at) RETURNING id ", &product)
	if err != nil {
		return fmt.Errorf("неудалось создать продкт %w", err)
	}
	return nil
}

func (p productRepo) ListProducts(limit, offset int) (*[]entity.Product, error) {
	var products []entity.Product
	err := p.db.Select(&products,
		"SELECT * FROM product ORDER BY id LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить продукты: %w", err)
	}
	return &products, nil
}
func (p productRepo) UpdateProduct(id uint, update *entity.Product) (*entity.Product, error) {
	params := map[string]interface{}{
		"id":          id,
		"name":        update.Name,
		"price":       update.Price,
		"description": update.Description,
		"stock":       update.Stock,
		"image_url":   update.ImageURL,
		"category_id": update.CategoryID,
	}
	_, err := p.db.NamedExec("UPDATE  product SET name = :name,price = :price,description = :description,stock = :stock,image_url = :image_url,category_id = :category_id,updated_at = NOW() WHERE id = :id", params)
	if err != nil {
		return nil, fmt.Errorf("не удалось оюновить продукт %w", err)
	}
	return update, err
}
func (p productRepo) DeleteProduct(id uint) error {
	_, err := p.db.Exec("DELETE FROM product WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("неудалось удалить продукт %w", err)
	}
	return err
}

func NewProdctRepo(db *sqlx.DB) *productRepo {
	return &productRepo{db: db}
}

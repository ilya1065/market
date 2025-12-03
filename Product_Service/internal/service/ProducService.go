package service

import "Product_Service/internal/entity"

type productRepo interface {
	Create(u *entity.Product) error
	GetProductByID(id uint) (*entity.Product, error)
	ListProducts(limit, offset int) (*[]entity.Product, error)
	UpdateProduct(id uint, update *entity.Product) (*entity.Product, error)
	DeleteProduct(id uint) error
}

func NewProductService(repo productRepo) *ProductService {
	return &ProductService{repo: repo}
}

type ProductService struct {
	repo productRepo
}

func (s ProductService) Create(u *entity.Product) error {
	err := s.repo.Create(u)
	return err
}

func (s ProductService) GetProductByID(id uint) (*entity.Product, error) {
	product, err := s.repo.GetProductByID(id)
	return product, err
}

func (s ProductService) ListProducts(limit, offset int) (*[]entity.Product, error) {
	products, err := s.repo.ListProducts(limit, offset)
	return products, err
}

func (s ProductService) UpdateProduct(id uint, update *entity.Product) (*entity.Product, error) {
	product, err := s.repo.UpdateProduct(id, update)
	return product, err
}

func (s ProductService) DeleteProduct(id uint) error {
	err := s.repo.DeleteProduct(id)
	return err
}

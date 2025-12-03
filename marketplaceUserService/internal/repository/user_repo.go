package repository

import (
	"gorm.io/gorm"
	"marketplace/internal/domain"
)

type UserRepository interface {
	Create(u *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	FindByID(id uint) (*domain.User, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(u *domain.User) error {
	return r.db.Create(u).Error
}

func (r *userRepo) FindByEmail(emsil string) (*domain.User, error) {
	var u domain.User
	if err := r.db.Where("email = ?", emsil).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) FindByID(id uint) (*domain.User, error) {
	var u domain.User
	if err := r.db.Find(id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

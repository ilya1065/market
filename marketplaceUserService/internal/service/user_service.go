package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"marketplace/internal/auth"
	"marketplace/internal/domain"
	"marketplace/internal/repository"
)

var (
	ErrEmailTaken   = errors.New("электронная почта уже используется ")
	ErrInvalidLogin = errors.New("неверный пароль или электронная почта")
	ErrTokenCreate  = errors.New("Ошибка при создании токена")
)

type UserService struct {
	repo repository.UserRepository
	Jwt  *auth.JWT
}

type Token struct {
	Refresh string
	Access  string
}

func NewUserService(repo repository.UserRepository, jwt *auth.JWT) *UserService {
	return &UserService{repo: repo, Jwt: jwt}
}

func (s *UserService) genTokens(u *domain.User) (*Token, error) {
	var t Token
	var err error
	t.Refresh, err = s.Jwt.RefreshGenerate(u.ID)
	if err != nil {
		return nil, ErrTokenCreate
	}
	t.Access, err = s.Jwt.AccessGenerate(u.ID)
	if err != nil {
		return nil, ErrTokenCreate
	}
	return &t, nil
}

func (s *UserService) Register(email, password string, hashFunc func(string) (string, error)) (*domain.User, *Token, error) {
	// Проверяем, что почта не занята
	if existing, _ := s.repo.FindByEmail(email); existing != nil {
		return nil, nil, ErrEmailTaken
	}

	hash, err := hashFunc(password)
	if err != nil {
		return nil, nil, err
	}

	u := &domain.User{
		Email:        email,
		PasswordHash: hash,
	}

	if err = s.repo.Create(u); err != nil {
		return nil, nil, err
	}

	t, err := s.genTokens(u)
	if err != nil {
		return nil, nil, err
	}

	return u, t, nil
}

func (s *UserService) Login(email, password string, check func(string, string) error) (*Token, error) {
	u, err := s.repo.FindByEmail(email)
	if err != nil || u == nil {
		return nil, ErrInvalidLogin
	}
	if err = check(u.PasswordHash, password); err != nil {
		return nil, ErrInvalidLogin
	}
	t, err := s.genTokens(u)
	if err != nil {
		return nil, err
	}

	return t, nil
}
func (s *UserService) GetByID(id uint) (*domain.User, error) {
	return s.repo.FindByID(id)
}
func (s *UserService) JWTParse(token string) (*jwt.RegisteredClaims, error) {
	return s.Jwt.Parse(token)
}

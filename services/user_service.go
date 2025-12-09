package services

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/nabilulilalbab/rental-mobil/models"
	"github.com/nabilulilalbab/rental-mobil/repositories"
)

type UserService interface {
	Register(user models.User) (models.User, error)
	Login(email, password string) (models.User, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) Register(user models.User) (models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}
	user.Password = string(hashedPassword)
	user.Role = "user" // Default role
	return s.repo.Create(user)
}

func (s *userService) Login(email, password string) (models.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return user, errors.New("email tidak ditemukan")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, errors.New("password salah")
	}

	return user, nil
}

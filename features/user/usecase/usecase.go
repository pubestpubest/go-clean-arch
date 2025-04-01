package usecase

import (
	"order-management/domain"
	"order-management/entity"
)

type userUsecase struct {
	repo domain.UserRepository
}

func NewUserUsecase(userRepository domain.UserRepository) domain.UserUsecase {
	return &userUsecase{repo: userRepository}
}

func (u *userUsecase) CreateUser(user entity.User) error {
	return u.repo.CreateUser(user)
}

package domain

import "order-management/entity"

type UserUsecase interface {
	CreateUser(user entity.User) error
	GetUserByID(id uint32) (entity.User, error)
	UpdateUser(user entity.User) error
}

type UserRepository interface {
	CreateUser(user entity.User) error
	GetUserByID(id uint32) (entity.User, error)
	UpdateUser(user entity.User) error
}

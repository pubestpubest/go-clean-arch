package domain

import "order-management/entity"

type UserUsecase interface {
	CreateUser(user entity.User) error
	UpdateUser(user entity.UserWithOutPassword) error
	Login(email string, password string) (string, error)
	GetUserByID(id uint32) (entity.UserWithOutPassword, error)
}

type UserRepository interface {
	CreateUser(user entity.User) error
	GetUserByID(id uint32) (entity.UserWithOutPassword, error)
	GetUserWithPasswordByEmail(email string) (entity.User, error)
	UpdateUser(user entity.UserWithOutPassword) error
	GetUserByEmail(email string) (entity.UserWithOutPassword, error)
}

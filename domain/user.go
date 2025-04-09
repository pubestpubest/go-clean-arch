package domain

import "order-management/entity"

type UserUsecase interface {
	CreateUser(user entity.User) error
	GetUserByID(id uint32) (entity.User, error)
	UpdateUser(user entity.User) error
	Login(email string, password string) (entity.User, error)
}

type UserRepository interface {
	CreateUser(user entity.User) error
	GetUserByID(id uint32) (entity.UserWithOutPassword, error)
	GetUserWithPasswordByEmail(email string) (entity.User, error)
	UpdateUser(user entity.User) error
	GetUserByEmail(email string) (entity.UserWithOutPassword, error)
}

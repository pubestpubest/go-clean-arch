package repository

import (
	"order-management/domain"
	"order-management/entity"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user entity.User) error {
	return r.db.Create(&user).Error
}

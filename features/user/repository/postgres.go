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

func (r *userRepository) GetUserByID(id uint32) (user entity.User, err error) {
	err = r.db.First(&user, id).Error
	return
}

func (r *userRepository) GetUserByEmail(email string) (user entity.User, err error) {
	err = r.db.Where("email = ?", email).First(&user).Error
	return
}

func (r *userRepository) UpdateUser(user entity.User) error {
	return r.db.Save(&user).Error
}

func (r *userRepository) DeleteUser(user entity.User) error {
	return r.db.Delete(&user).Error
}

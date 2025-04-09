package repository

import (
	"order-management/domain"
	"order-management/entity"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user entity.User) error {
	if err := r.db.Create(&user).Error; err != nil {
		err = errors.Wrap(err, "[UserRepository.CreateUser]: failed to create user")
		return err
	}
	return nil
}

func (r *userRepository) GetUserByID(id uint32) (user entity.UserWithOutPassword, err error) {
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("[UserRepository.GetUserByID]: user not found")
			return entity.UserWithOutPassword{}, err
		}
		err = errors.Wrap(err, "[UserRepository.GetUserByID]: failed to get user by id")
		return entity.UserWithOutPassword{}, err
	}
	return user, nil
}

func (r *userRepository) GetUserByEmail(email string) (user entity.UserWithOutPassword, err error) {
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("[UserRepository.GetUserByEmail]: user not found")
			return entity.UserWithOutPassword{}, err
		}
		err = errors.Wrap(err, "[UserRepository.GetUserByEmail]: failed to get user by email")
		return entity.UserWithOutPassword{}, err
	}
	return user, nil
}

func (r *userRepository) GetUserWithPasswordByEmail(email string) (user entity.User, err error) {
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("[UserRepository.GetUserWithPasswordByEmail]: user not found")
			return entity.User{}, err
		}
		err = errors.Wrap(err, "[UserRepository.GetUserWithPasswordByEmail]: failed to get user with password by email")
		return entity.User{}, err
	}
	return user, nil
}

func (r *userRepository) UpdateUser(user entity.User) error {
	return r.db.Save(&user).Error
}

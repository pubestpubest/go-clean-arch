package usecase

import (
	"order-management/domain"
	"order-management/entity"
	"order-management/utils"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	repo domain.UserRepository
}

func NewUserUsecase(userRepository domain.UserRepository) domain.UserUsecase {
	return &userUsecase{repo: userRepository}
}

func (u *userUsecase) CreateUser(user entity.User) error {
	log.Trace("Entering function CreateUser()")
	defer log.Trace("Exiting function CreateUser()")

	log.WithFields(log.Fields{
		"user": user,
	}).Debug("Creating user")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		err = errors.Wrap(err, "[UserUsecase.CreateUser]: failed to hash password")
		return err
	}
	user.Password = string(hashedPassword)

	if err := u.repo.CreateUser(user); err != nil {
		if err.Error() == "[UserRepository.CreateUser]: user already exists" {
			err = errors.New("[UserUsecase.CreateUser]: user already exists")
			return err
		}
		err = errors.Wrap(err, "[UserUsecase.CreateUser]: failed to create user")
		return err
	}
	return nil
}

func (u *userUsecase) UpdateUser(user entity.UserWithOutPassword) error {
	log.Trace("Entering function UpdateUser()")
	defer log.Trace("Exiting function UpdateUser()")

	log.WithFields(log.Fields{
		"user": user,
	}).Debug("Updating user")

	if err := u.repo.UpdateUser(user); err != nil {
		err = errors.Wrap(err, "[UserUsecase.UpdateUser]: failed to update user")
		return err
	}
	return nil
}

func (u *userUsecase) Login(email string, password string) (string, error) {
	log.Trace("Entering function Login()")
	defer log.Trace("Exiting function Login()")

	log.WithFields(log.Fields{
		"email":    email,
		"password": password,
	}).Debug("Logging in user")

	credentials, err := u.repo.GetUserWithPasswordByEmail(email)
	if err != nil {
		if err.Error() == "[UserRepository.GetUserWithPasswordByEmail]: user not found" {
			err = errors.New("[UserUsecase.Login]: user not found")
			return "", err
		}
		err = errors.Wrap(err, "[UserUsecase.Login]: failed to get user with password by email")
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(credentials.Password), []byte(password)); err != nil {
		err = errors.New("[UserUsecase.Login]: invalid password")
		return "", err
	}

	t, err := utils.GenerateJWT(map[string]interface{}{
		"id":      credentials.ID,
		"email":   credentials.Email,
		"address": credentials.Address,
	}, []byte(viper.GetString("jwt.usersecret")))

	if err != nil {
		err = errors.Wrap(err, "[UserUsecase.Login]: failed to generate user jwt")
		return "", err
	}

	return t, nil
}
func (u *userUsecase) GetUserByID(id uint32) (entity.UserWithOutPassword, error) {
	log.Trace("Entering function GetUserByID()")
	defer log.Trace("Exiting function GetUserByID()")

	log.WithFields(log.Fields{
		"id": id,
	}).Debug("Getting user by id")

	user, err := u.repo.GetUserByID(id)
	if err != nil {
		err = errors.Wrap(err, "[UserUsecase.GetUserByID]: failed to get user by id")
		return entity.UserWithOutPassword{}, err
	}
	return user, nil
}

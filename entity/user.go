package entity

import "github.com/golang-jwt/jwt/v5"

type User struct {
	ID       uint32 `gorm:"primary_key"`
	Email    string `gorm:"unique;not null"`
	Address  string
	Password string  `gorm:"not null"`
	Orders   []Order `gorm:"foreignKey:UserID"`
}

type UserWithOutPassword struct {
	ID      uint32 `json:"id"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

type UserJWT struct {
	ID      uint32 `json:"id"`
	Email   string `json:"email"`
	Address string `json:"address"`
	jwt.RegisteredClaims
}

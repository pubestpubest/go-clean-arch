package entity

import "github.com/golang-jwt/jwt/v5"

type Shop struct {
	ID          uint32 `gorm:"primary_key"`
	Name        string `gorm:"not null;unique"`
	Description string
	Password    string    `gorm:"not null"`
	Products    []Product `gorm:"foreignKey:ShopID"`
}

type ShopResponse struct {
	ID          uint32 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	jwt.RegisteredClaims
}

type ShopWithProducts struct {
	ID          uint32            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Products    []ProductResponse `json:"products"`
}

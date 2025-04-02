package response

import "github.com/golang-jwt/jwt/v5"

type Shop struct {
	ID          uint32 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	jwt.RegisteredClaims
}
type ShopWithProducts struct {
	ID          uint32    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Products    []Product `json:"products"`
}

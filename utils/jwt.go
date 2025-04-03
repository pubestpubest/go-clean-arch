package utils

import (
	"errors"
	"order-management/entity"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateShopJWT(claims *entity.ShopResponse) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(ViperGetString("jwt.secret")))
}

func ValidateShopJWT(tokenString string) (*entity.ShopResponse, error) {
	token, err := jwt.ParseWithClaims(tokenString, &entity.ShopResponse{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(ViperGetString("jwt.secret")), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*entity.ShopResponse)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

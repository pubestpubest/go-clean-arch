package utils

import (
	"order-management/entity"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

func GenerateShopJWT(claims *entity.ShopResponse) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(ViperGetString("jwt.secret")))
	if err != nil {
		return "", errors.Wrap(err, "[utils.GenerateShopJWT]: failed to generate shop jwt")
	}
	return tokenString, nil
}

func ValidateShopJWT(tokenString string) (*entity.ShopResponse, error) {
	token, err := jwt.ParseWithClaims(tokenString, &entity.ShopResponse{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(ViperGetString("jwt.secret")), nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "[utils.ValidateShopJWT]: failed to parse shop jwt")
	}
	claims, ok := token.Claims.(*entity.ShopResponse)
	if !ok || !token.Valid {
		return nil, errors.Wrap(err, "[utils.ValidateShopJWT]: invalid token")
	}
	return claims, nil
}

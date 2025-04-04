package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

func GenerateJWT(payload map[string]interface{}, secret []byte) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	for k, v := range payload {
		claims[k] = v
	}
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", errors.Wrap(err, "[utils.GenerateShopJWT]: failed to generate shop jwt")
	}
	return tokenString, nil
}

func ValidateJWT(tokenString string, secret []byte) (*jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "[utils.ValidateShopJWT]: failed to parse shop jwt")
	}
	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.Wrap(err, "[utils.ValidateShopJWT]: invalid token")
	}
	return claims, nil
}

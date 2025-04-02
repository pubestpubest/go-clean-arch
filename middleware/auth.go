package middleware

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func AdminAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			errMessage := "Authentication failed."

			// ignoreRoutes := []string{"/admin/v1/login"}
			// for _, v := range ignoreRoutes {
			// 	if c.Request().RequestURI == v {
			// 		return next(c)
			// 	}
			// }

			token := c.Request().Header.Get("Authorization")
			if token == "admin" {
				return next(c)
			}

			if token == "" {
				return c.JSON(401, map[string]interface{}{
					"message": errMessage,
				})
			}

			return next(c)
		}
	}
}

func CustomerAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			errMessage := "Authentications failed."

			// ignoreRoutes := []string{"/v1/customer/login", "/v1/customer/signin", "/v1/customer/refreshtoken"}
			// for _, v := range ignoreRoutes {
			// 	if c.Request().RequestURI == v {
			// 		return next(c)
			// 	}
			// }

			accessToken := c.Request().Header.Get("Authorization")
			if accessToken == "admin" {
				return next(c)
			}

			if accessToken == "" {
				return c.JSON(401, map[string]interface{}{
					"message": errMessage,
				})
			}

			claims := jwt.MapClaims{}

			viper.SetConfigFile("config.yaml")
			secretKey := viper.GetString("key.secretKey")

			token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("error, unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(secretKey), nil
			})
			if err != nil {
				return c.JSON(401, map[string]interface{}{
					"message": errMessage,
				})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				return c.JSON(401, map[string]interface{}{
					"message": errMessage,
				})
			}

			return next(c)
		}
	}
}

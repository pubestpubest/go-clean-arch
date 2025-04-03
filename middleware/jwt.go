package middleware

import (
	"net/http"
	"order-management/entity"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func ShopAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("token")
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
			}
			token, err := jwt.ParseWithClaims(cookie.Value, &entity.ShopResponse{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(viper.GetString("jwt.secret")), nil
			})
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
			}
			claims, ok := token.Claims.(*entity.ShopResponse)
			if !ok || !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
			}
			c.Set("shop", claims)
			return next(c)
		}
	}
}

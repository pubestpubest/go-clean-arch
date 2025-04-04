package middleware

import (
	"net/http"
	"order-management/entity"
	"order-management/utils"
	"strings"

	"github.com/labstack/echo/v4"
)

func ShopAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			bearerToken := c.Request().Header.Get("Authorization")
			if bearerToken == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, entity.ResponseError{
					Error: "No authorization header found",
				})
			}
			str := strings.Split(bearerToken, " ")
			if len(str) != 2 {
				return echo.NewHTTPError(http.StatusUnauthorized, entity.ResponseError{
					Error: "Invalid authorization header",
				})
			}
			token := str[1]
			claims, err := utils.ValidateShopJWT(token)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, entity.ResponseError{
					Error: err.Error(),
				})
			}
			c.Set("shop", claims)
			return next(c)
		}
	}
}

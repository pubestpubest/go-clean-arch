package middleware

import (
	"net/http"
	"order-management/entity"
	"order-management/utils"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
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
			claims, err := utils.ValidateJWT(token, []byte(viper.GetString("jwt.shopsecret")))
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, entity.ResponseError{
					Error: err.Error(),
				})
			}

			// Convert MapClaims to ShopWithOutPassword
			shopClaims := &entity.ShopWithOutPassword{
				ID:          uint32((*claims)["id"].(float64)),
				Name:        (*claims)["name"].(string),
				Description: (*claims)["description"].(string),
			}

			c.Set("shop", shopClaims)

			return next(c)
		}
	}
}

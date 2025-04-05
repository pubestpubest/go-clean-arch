package middleware

import (
	"errors"
	"net/http"
	"order-management/entity"
	"order-management/utils"
	"strings"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Todo: only failed case
// Happy case is not that important
func ShopAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			bearerToken := c.Request().Header.Get("Authorization")
			if bearerToken == "" {
				err := errors.New("[Middleware.ShopAuth]: no authorization header found")
				log.WithError(err).Warn("Missing authorization header")
				return echo.NewHTTPError(http.StatusUnauthorized, entity.ResponseError{
					Error: utils.StandardError(err),
				})
			}

			str := strings.Split(bearerToken, " ")
			if len(str) != 2 {
				err := errors.New("[Middleware.ShopAuth]: invalid authorization header format")
				log.WithError(err).Warn("Invalid authorization header format")
				return echo.NewHTTPError(http.StatusUnauthorized, entity.ResponseError{
					Error: utils.StandardError(err),
				})
			}

			token := str[1]
			claims, err := utils.ValidateJWT(token, []byte(viper.GetString("jwt.shopsecret")))
			if err != nil {
				err := errors.New("[Middleware.ShopAuth]: JWT validation failed")
				log.WithError(err).Error("JWT validation failed")
				return echo.NewHTTPError(http.StatusUnauthorized, entity.ResponseError{
					Error: utils.StandardError(err),
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

func UserAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
		}
	}
}

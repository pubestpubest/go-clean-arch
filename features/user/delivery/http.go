package delivery

import (
	"net/http"
	"order-management/domain"
	"order-management/entity"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type Handler struct {
	usecase domain.UserUsecase
}

func NewHandler(e *echo.Group, u domain.UserUsecase) *Handler {
	h := Handler{usecase: u}

	e.POST("/users/register", h.CreateUser)
	e.POST("/users/login", h.Login)
	e.GET("/users/:id", h.GetUserByID)
	e.PUT("/users/:id", h.UpdateUser)
	e.GET("/users/me", h.ReadToken)

	return &h
}

func (h *Handler) CreateUser(c echo.Context) error {
	req := entity.User{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, entity.ResponseError{Error: err.Error()})
	}
	if req.Password == "" {
		return c.JSON(http.StatusBadRequest, entity.ResponseError{Error: "Password is required"})
	}

	if err := h.usecase.CreateUser(req); err != nil {
		return c.JSON(http.StatusInternalServerError, entity.ResponseError{Error: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) GetUserByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, entity.ResponseError{Error: err.Error()})
	}

	user, err := h.usecase.GetUserByID(uint32(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, entity.ResponseError{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *Handler) UpdateUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, entity.ResponseError{Error: err.Error()})
	}

	user, err := h.usecase.GetUserByID(uint32(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, entity.ResponseError{Error: err.Error()})
	}

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, entity.ResponseError{Error: err.Error()})
	}

	if err := h.usecase.UpdateUser(user); err != nil {
		return c.JSON(http.StatusInternalServerError, entity.ResponseError{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, user)
}

func (h *Handler) Login(c echo.Context) error {
	req := entity.User{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, entity.ResponseError{Error: err.Error()})
	}

	user, err := h.usecase.Login(req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, entity.ResponseError{Error: err.Error()})
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"id":    user.ID,
	})
	t, err := token.SignedString([]byte(viper.GetString("jwt.secret")))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, entity.ResponseError{Error: err.Error()})
	}
	cookie := &http.Cookie{
		Name:  "token",
		Value: t,
		Path:  "/",
	}
	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, user)
}

func (h *Handler) ReadToken(c echo.Context) error {
	cookie, err := c.Cookie("token")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, entity.ResponseError{Error: "Unauthorized"})
	}
	token, err := jwt.ParseWithClaims(cookie.Value, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("jwt.secret")), nil
	})
	if err != nil {
		return c.JSON(http.StatusUnauthorized, entity.ResponseError{Error: "Unauthorized"})
	}
	claims := token.Claims.(jwt.MapClaims)
	address, err := h.usecase.GetAddressByUserID(uint32(claims["id"].(float64)))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, entity.ResponseError{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"address": address,
		"email":   claims["email"],
		"id":      claims["id"],
	})
}

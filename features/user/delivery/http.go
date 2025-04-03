package delivery

import (
	"net/http"
	"order-management/domain"
	"order-management/entity"
	"strconv"

	"order-management/middleware"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	usecase domain.UserUsecase
}

func NewHandler(e *echo.Group, u domain.UserUsecase) *Handler {
	h := Handler{usecase: u}

	e.POST("/users", h.CreateUser)
	e.GET("/users/:id", h.GetUserByID) //Should be used only by Shop
	e.PUT("/users/:id", h.UpdateUser)

	publicGroup := e.Group("")
	publicGroup.POST("/users/register", h.CreateUser)
	publicGroup.POST("/users/login", h.Login)

	authGroup := e.Group("")
	authGroup.Use(middleware.UserAuth())

	return &h
}

func (h *Handler) Login(c echo.Context) error {
	//Bind
	req := entity.User{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, entity.ResponseError{Error: err.Error()})
	}
	//Check if email and password are provided
	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, entity.ResponseError{Error: "email and password are required"})
	}
	//Login
	user, err := h.usecase.Login(req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, entity.ResponseError{Error: err.Error()})
	}
	//Check if password is correct
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, entity.ResponseError{Error: "invalid password"})
	}
	//Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
	})
	t, err := token.SignedString([]byte(viper.GetString("jwt.secret")))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, entity.ResponseError{Error: err.Error()})
	}
	//Set cookie
	cookie := &http.Cookie{
		Name:  "token",
		Value: t,
		Path:  "/",
	}
	c.SetCookie(cookie)
	//Return user
	return c.JSON(http.StatusOK, user)
}

func (h *Handler) CreateUser(c echo.Context) error {
	req := entity.User{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, entity.ResponseError{Error: err.Error()})
	}
	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, entity.ResponseError{Error: "email and password are required"})
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

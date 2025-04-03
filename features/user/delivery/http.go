package delivery

import (
	"net/http"
	"order-management/domain"
	"order-management/entity"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	usecase domain.UserUsecase
}

func NewHandler(e *echo.Group, u domain.UserUsecase) *Handler {
	h := Handler{usecase: u}

	e.POST("/users", h.CreateUser)
	e.GET("/users/:id", h.GetUserByID)
	e.PUT("/users/:id", h.UpdateUser)

	return &h
}

func (h *Handler) CreateUser(c echo.Context) error {
	req := entity.User{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, entity.ResponseError{Error: err.Error()})
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

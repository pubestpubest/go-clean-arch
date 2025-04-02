package delivery

import (
	"net/http"
	"order-management/domain"
	"order-management/entity"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	usecase domain.UserUsecase
}

func NewHandler(e *echo.Group, u domain.UserUsecase) *Handler {
	h := Handler{usecase: u}

	e.POST("/users", h.CreateUser)

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

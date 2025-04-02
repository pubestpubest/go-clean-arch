package delivery

import (
	"net/http"
	"order-management/domain"
	"order-management/entity"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	usecase domain.ShopUsecase
}

func NewHandler(e *echo.Group, u domain.ShopUsecase) *Handler {
	h := Handler{usecase: u}

	e.POST("/:shop_id/products", h.CreateProduct)
	e.GET("/shops", h.GetAllShops)
	e.POST("/shops", h.CreateShop)
	return &h
}

func (h *Handler) CreateProduct(c echo.Context) error {
	req := entity.Product{}
	shopID, err := strconv.ParseUint(c.Param("shop_id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := h.usecase.CreateProduct(req, uint32(shopID)); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) GetAllShops(c echo.Context) error {
	shops, err := h.usecase.GetAllShopsWithProducts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, shops)
}

func (h *Handler) CreateShop(c echo.Context) error {
	req := entity.Shop{}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := h.usecase.CreateShop(req); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

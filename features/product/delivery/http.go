package delivery

import (
	"net/http"
	"order-management/domain"
	"order-management/entity"
	"order-management/utils"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	usecase domain.ProductUsecase
}

func NewHandler(e *echo.Group, u domain.ProductUsecase) *Handler {
	h := Handler{usecase: u}

	publicGroup := e.Group("")
	publicGroup.GET("", h.GetAllProducts)
	publicGroup.GET("/:productID", h.GetProductByID)
	return &h
}

func (h *Handler) GetProductByID(c echo.Context) error {
	productID := c.Param("productID")
	productIDUint, err := strconv.ParseUint(productID, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, entity.ResponseError{
			Error: utils.StandardError(err),
		})
	}
	product, err := h.usecase.GetProductByID(uint32(productIDUint))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, entity.ResponseError{
			Error: utils.StandardError(err),
		})
	}
	return c.JSON(http.StatusOK, entity.Response{
		Success: true,
		Message: "Product fetched successfully",
		Data:    product,
		Status:  http.StatusOK,
	})
}
func (h *Handler) GetAllProducts(c echo.Context) error {
	products, err := h.usecase.GetAllProducts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, entity.ResponseError{
			Error: utils.StandardError(err),
		})
	}
	return c.JSON(http.StatusOK, entity.Response{
		Success: true,
		Message: "Products fetched successfully",
		Data:    products,
		Status:  http.StatusOK,
	})
}

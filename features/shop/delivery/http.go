package delivery

import (
	"net/http"
	"order-management/domain"
	"order-management/entity"
	"order-management/utils"
	"strconv"

	"order-management/middleware"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Handler struct {
	usecase domain.ShopUsecase
}

func NewHandler(e *echo.Group, u domain.ShopUsecase) *Handler {
	h := Handler{usecase: u}
	// Public group - no authentication required
	publicGroup := e.Group("")
	publicGroup.GET("", h.GetAllShops)                           // Anyone can view shops
	publicGroup.POST("/register", h.CreateShop)                  // Public registration
	publicGroup.POST("/login", h.Login)                          // Public login
	publicGroup.GET("/:shop_id/products", h.GetProductsByShopID) // Anyone can view products

	// Authenticated group - requires JWT
	authGroup := e.Group("")
	authGroup.Use(middleware.ShopAuth())
	authGroup.POST("/products", h.CreateProduct)               // Only shop owner can create products
	authGroup.PUT("/products/:product_id", h.UpdateProduct)    // Only shop owner can update their products
	authGroup.DELETE("/products/:product_id", h.DeleteProduct) // Only shop owner can delete their products
	authGroup.GET("/me", h.ReadToken)                          // Get current shop profile from JWT
	authGroup.POST("/logout", h.Logout)                        // Logout requires JWT
	authGroup.GET("/profile", h.GetShopProfile)                // Get detailed profile requires JWT

	return &h
}

func (h *Handler) GetShopProfile(c echo.Context) error {
	log.Trace("Entering function GetShopProfile()")
	shopClaims, ok := c.Get("shop").(*entity.ShopWithOutPassword)
	if !ok {
		err := errors.New("[Handler.GetShopProfile]: no shop claims found")
		log.WithError(err).Error("Failed to get shop claims from context")
		return c.JSON(http.StatusUnauthorized, entity.ResponseError{
			Error: utils.StandardError(err),
		})
	}

	log.WithField("shopName", shopClaims.Name).Debug("Attempting to retrieve shop profile")
	shop, err := h.usecase.GetShopByName(shopClaims.Name)
	if err != nil {
		if err.Error() == "[ShopUsecase.GetShopByName]: shop not found" {
			// If this happens, it means the shop name is not in the database
			// Or the JWT secret is compromised
			log.WithFields(log.Fields{
				"shopName": shopClaims.Name,
				"error":    err,
			}).Warn("Shop not found")
			return c.JSON(http.StatusNotFound, entity.ResponseError{
				Error: utils.StandardError(errors.Wrap(err, "[Handler.GetShopProfile]: shop not found")),
			})
		}
		log.WithFields(log.Fields{
			"shopName": shopClaims.Name,
			"error":    err,
		}).Error("Internal server error while retrieving shop profile")
		return c.JSON(http.StatusInternalServerError, entity.ResponseError{
			Error: utils.StandardError(errors.Wrap(err, "[Handler.GetShopProfile]: internal server error")),
		})
	}

	log.WithField("shopName", shopClaims.Name).Info("Shop profile retrieved successfully")
	return c.JSON(http.StatusOK, entity.Response{
		Success: true,
		Message: "Shop profile retrieved successfully",
		Data:    shop,
		Status:  http.StatusOK,
	})
}

func (h *Handler) GetProductsByShopID(c echo.Context) error {
	shopID, err := strconv.ParseUint(c.Param("shop_id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, entity.ResponseError{
			Error: errors.Wrap(err, "[Handler.GetProductsByShopID]: invalid shop id").Error(),
		})
	}

	products, err := h.usecase.GetProductsByShopID(uint32(shopID))
	if err != nil {
		if err.Error() == "[ShopUsecase.GetProductsByShopID]: shop not found" {
			return c.JSON(http.StatusNotFound, entity.ResponseError{
				Error: errors.Wrap(err, "[Handler.GetProductsByShopID]: shop not found").Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, entity.ResponseError{
			Error: errors.Wrap(err, "[Handler.GetProductsByShopID]: internal server error").Error(),
		})
	}

	return c.JSON(http.StatusOK, entity.Response{
		Success: true,
		Message: "Products retrieved successfully",
		Data:    products,
		Status:  http.StatusOK,
	})
}

func (h *Handler) DeleteProduct(c echo.Context) error {
	productID, err := strconv.ParseUint(c.Param("product_id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, entity.ResponseError{
			Error: errors.Wrap(err, "[Handler.DeleteProduct]: invalid product id").Error(),
		})
	}

	shop, ok := c.Get("shop").(*entity.ShopWithOutPassword)
	if !ok {
		return c.JSON(http.StatusUnauthorized, entity.ResponseError{
			Error: "[Handler.DeleteProduct]: no shop claims found",
		})
	}

	req := entity.ProductManagementRequest{
		ShopID:    shop.ID,
		ProductID: uint32(productID),
	}

	if err := h.usecase.DeleteProduct(&req); err != nil {
		switch err.Error() {
		case "[ShopUsecase.DeleteProduct]: shop not found":
			return c.JSON(http.StatusNotFound, entity.ResponseError{
				Error: errors.Wrap(err, "[Handler.DeleteProduct]: shop not found").Error(),
			})
		case "[ShopUsecase.DeleteProduct]: product not found":
			return c.JSON(http.StatusNotFound, entity.ResponseError{
				Error: errors.Wrap(err, "[Handler.DeleteProduct]: product not found").Error(),
			})
		case "[ShopUsecase.DeleteProduct]: product does not belong to shop":
			return c.JSON(http.StatusForbidden, entity.ResponseError{
				Error: errors.Wrap(err, "[Handler.DeleteProduct]: product does not belong to shop").Error(),
			})
		default:
			return c.JSON(http.StatusInternalServerError, entity.ResponseError{
				Error: errors.Wrap(err, "[Handler.DeleteProduct]: internal server error").Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, entity.Response{
		Success: true,
		Message: "Product deleted successfully",
		Status:  http.StatusOK,
	})
}

func (h *Handler) UpdateProduct(c echo.Context) error {
	productID, err := strconv.ParseUint(c.Param("product_id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, entity.ResponseError{
			Error: errors.Wrap(err, "[Handler.UpdateProduct]: invalid product id").Error(),
		})
	}

	shop, ok := c.Get("shop").(*entity.ShopWithOutPassword)
	if !ok {
		return c.JSON(http.StatusUnauthorized, entity.ResponseError{
			Error: "[Handler.UpdateProduct]: no shop claims found",
		})
	}

	product := entity.Product{}
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, entity.ResponseError{
			Error: errors.Wrap(err, "[Handler.UpdateProduct]: invalid product").Error(),
		})
	}

	req := entity.ProductManagementRequest{
		ShopID:    shop.ID,
		ProductID: uint32(productID),
	}

	if err := h.usecase.UpdateProduct(&req, &product); err != nil {
		switch err.Error() {
		case "[ShopUsecase.UpdateProduct]: shop not found":
			return c.JSON(http.StatusNotFound, entity.ResponseError{
				Error: errors.Wrap(err, "[Handler.UpdateProduct]: shop not found").Error(),
			})
		case "[ShopUsecase.UpdateProduct]: product not found":
			return c.JSON(http.StatusNotFound, entity.ResponseError{
				Error: errors.Wrap(err, "[Handler.UpdateProduct]: product not found").Error(),
			})
		case "[ShopUsecase.UpdateProduct]: product does not belong to shop":
			return c.JSON(http.StatusForbidden, entity.ResponseError{
				Error: errors.Wrap(err, "[Handler.UpdateProduct]: product does not belong to shop").Error(),
			})
		default:
			return c.JSON(http.StatusInternalServerError, entity.ResponseError{
				Error: errors.Wrap(err, "[Handler.UpdateProduct]: internal server error").Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, entity.Response{
		Success: true,
		Message: "Product updated successfully",
		Status:  http.StatusOK,
	})
}

func (h *Handler) Logout(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func (h *Handler) CreateProduct(c echo.Context) error {
	shopClaims, ok := c.Get("shop").(*entity.ShopWithOutPassword)
	if !ok {
		return c.JSON(http.StatusUnauthorized, entity.ResponseError{
			Error: "No shop claims found",
		})
	}

	req := entity.Product{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, entity.ResponseError{
			Error: errors.Wrap(err, "[Handler.CreateProduct]: invalid product").Error(),
		})
	}

	if err := h.usecase.CreateProduct(req, shopClaims.ID); err != nil {
		return c.JSON(http.StatusInternalServerError, entity.ResponseError{
			Error: errors.Wrap(err, "[Handler.CreateProduct]: internal server error").Error(),
		})
	}

	return c.JSON(http.StatusOK, entity.Response{
		Success: true,
		Message: "Product created successfully",
		Status:  http.StatusOK,
	})
}

func (h *Handler) GetAllShops(c echo.Context) error {
	shops, err := h.usecase.GetAllShopsWithProducts()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, entity.ResponseError{
				Error: errors.Wrap(err, "[Handler.GetAllShops]: no shops found").Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, entity.ResponseError{
			Error: errors.Wrap(err, "[Handler.GetAllShops]: internal server error").Error(),
		})
	}
	return c.JSON(http.StatusOK, entity.Response{
		Success: true,
		Message: "Shops retrieved successfully",
		Data:    shops,
		Status:  http.StatusOK,
	})
}

func (h *Handler) CreateShop(c echo.Context) error {
	req := entity.Shop{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, entity.ResponseError{
			Error: errors.Wrap(err, "[Handler.CreateShop]: invalid shop").Error(),
		})
	}

	if req.Password == "" {
		return c.JSON(http.StatusBadRequest, entity.ResponseError{
			Error: "password is required",
		})
	}

	if err := h.usecase.CreateShop(req); err != nil {
		return c.JSON(http.StatusInternalServerError, entity.ResponseError{
			Error: errors.Wrap(err, "[Handler.CreateShop]: internal server error").Error(),
		})
	}

	return c.JSON(http.StatusOK, entity.Response{
		Success: true,
		Message: "Shop created successfully",
		Status:  http.StatusOK,
	})
}

func (h *Handler) Login(c echo.Context) error {
	req := entity.Shop{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, entity.ResponseError{
			Error: errors.Wrap(err, "[Handler.Login]: invalid shop").Error(),
		})
	}

	if req.Name == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, entity.ResponseError{
			Error: "[Handler.Login]: name and password are required",
		})
	}

	token, err := h.usecase.Login(req.Name, req.Password)
	if err != nil {
		if err.Error() == "[ShopUsecase.Login]: shop not found" {
			return c.JSON(http.StatusNotFound, entity.ResponseError{
				Error: errors.Wrap(err, "[Handler.Login]: shop not found").Error(),
			})
		}
		if err.Error() == "[ShopUsecase.Login]: invalid password" {
			return c.JSON(http.StatusUnauthorized, entity.ResponseError{
				Error: errors.Wrap(err, "[Handler.Login]: invalid password").Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, entity.ResponseError{
			Error: errors.Wrap(err, "[Handler.Login]: internal server error").Error(),
		})
	}

	c.Response().Header().Set("Authorization", "Bearer "+token)

	return c.JSON(http.StatusOK, entity.Response{
		Success: true,
		Message: "Login successful",
		Status:  http.StatusOK,
	})
}

func (h *Handler) ReadToken(c echo.Context) error {
	shopClaims, ok := c.Get("shop").(*entity.ShopWithOutPassword)
	if !ok {
		return c.JSON(http.StatusUnauthorized, entity.ResponseError{
			Error: "unauthorized",
		})
	}

	return c.JSON(http.StatusOK, shopClaims)
}

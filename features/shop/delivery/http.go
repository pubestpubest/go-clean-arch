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
	defer log.Trace("Exiting function GetShopProfile()")

	shopClaims, ok := c.Get("shop").(*entity.ShopJWT)

	log.WithFields(log.Fields{
		"shopClaims": shopClaims,
	}).Debug("Shop claims from context")

	if !ok {
		err := errors.New("[Handler.GetShopProfile]: no shop claims found")

		log.WithError(err).Warn("No shop claims found in context")

		return c.JSON(http.StatusUnauthorized, entity.ResponseError{
			Error: utils.StandardError(err),
		})
	}

	log.WithField("shopName", shopClaims.Name).Debug("Attempting to retrieve shop profile")

	shop, err := h.usecase.GetShopByName(shopClaims.Name)
	if err != nil {
		if err.Error() == "[ShopUsecase.GetShopByName]: shop not found" {
			// If this happens, it means the shop name is not in the database
			// the JWT secret is compromised or the shop is deleted
			err = errors.Wrap(err, "[Handler.GetShopProfile]: shop not found")

			log.WithFields(log.Fields{
				"shopName": shopClaims.Name,
			}).WithError(err).Warn("Shop not found")

			return c.JSON(http.StatusNotFound, entity.ResponseError{
				Error: utils.StandardError(err),
			})
		}
		err = errors.Wrap(err, "[Handler.GetShopProfile]: internal server error")

		log.WithFields(log.Fields{
			"shopName": shopClaims.Name,
		}).WithError(err).Error("Internal server error while retrieving shop profile")

		return c.JSON(http.StatusInternalServerError, entity.ResponseError{
			Error: utils.StandardError(err),
		})
	}

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
		err = errors.Wrap(err, "[Handler.GetProductsByShopID]: invalid shop id")

		log.WithError(err).Warn("Invalid shop ID format")

		return c.JSON(http.StatusBadRequest, entity.ResponseError{
			Error: err.Error(),
		})
	}

	products, err := h.usecase.GetProductsByShopID(uint32(shopID))
	if err != nil {
		if err.Error() == "[ShopUsecase.GetProductsByShopID]: shop not found" {
			err = errors.Wrap(err, "[Handler.GetProductsByShopID]: shop not found")

			log.WithFields(log.Fields{
				"shopID": shopID,
			}).WithError(err).Warn("Shop not found")

			return c.JSON(http.StatusNotFound, entity.ResponseError{
				Error: err.Error(),
			})
		}
		err = errors.Wrap(err, "[Handler.GetProductsByShopID]: internal server error")

		log.WithFields(log.Fields{
			"shopID": shopID,
		}).WithError(err).Error("Internal server error while getting products")

		return c.JSON(http.StatusInternalServerError, entity.ResponseError{
			Error: err.Error(),
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
		err = errors.Wrap(err, "[Handler.DeleteProduct]: invalid product id")

		log.WithError(err).Warn("Invalid product ID format")

		return c.JSON(http.StatusBadRequest, entity.ResponseError{
			Error: err.Error(),
		})
	}

	shop, ok := c.Get("shop").(*entity.ShopJWT)
	if !ok {
		err := errors.New("[Handler.DeleteProduct]: no shop claims found")

		log.Warn("No shop claims found in context")

		return c.JSON(http.StatusUnauthorized, entity.ResponseError{
			Error: err.Error(),
		})
	}

	req := entity.ProductManagementRequest{
		ShopID:    shop.ID,
		ProductID: uint32(productID),
	}

	if err := h.usecase.DeleteProduct(&req); err != nil {
		switch err.Error() {
		case "[ShopUsecase.DeleteProduct]: shop not found":
			err = errors.Wrap(err, "[Handler.DeleteProduct]: shop not found")

			log.WithFields(log.Fields{
				"shopID": shop.ID,
			}).WithError(err).Warn("Shop not found")

			return c.JSON(http.StatusNotFound, entity.ResponseError{
				Error: err.Error(),
			})
		case "[ShopUsecase.DeleteProduct]: product not found":
			err = errors.Wrap(err, "[Handler.DeleteProduct]: product not found")

			log.WithFields(log.Fields{
				"productID": productID,
			}).WithError(err).Warn("Product not found")

			return c.JSON(http.StatusNotFound, entity.ResponseError{
				Error: err.Error(),
			})
		case "[ShopUsecase.DeleteProduct]: product does not belong to shop":
			err = errors.Wrap(err, "[Handler.DeleteProduct]: product does not belong to shop")

			log.WithFields(log.Fields{
				"shopID":    shop.ID,
				"productID": productID,
			}).WithError(err).Warn("Product does not belong to shop")

			return c.JSON(http.StatusForbidden, entity.ResponseError{
				Error: err.Error(),
			})
		default:
			err = errors.Wrap(err, "[Handler.DeleteProduct]: internal server error")

			log.WithFields(log.Fields{
				"shopID":    shop.ID,
				"productID": productID,
			}).WithError(err).Error("Internal server error while deleting product")

			return c.JSON(http.StatusInternalServerError, entity.ResponseError{
				Error: err.Error(),
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
		err = errors.Wrap(err, "[Handler.UpdateProduct]: invalid product id")

		log.WithError(err).Warn("Invalid product ID format")

		return c.JSON(http.StatusBadRequest, entity.ResponseError{
			Error: err.Error(),
		})
	}

	shop, ok := c.Get("shop").(*entity.ShopJWT)
	if !ok {
		err := errors.New("[Handler.UpdateProduct]: no shop claims found")

		log.Warn("No shop claims found in context")

		return c.JSON(http.StatusUnauthorized, entity.ResponseError{
			Error: err.Error(),
		})
	}

	product := entity.Product{}
	if err := c.Bind(&product); err != nil {
		err = errors.Wrap(err, "[Handler.UpdateProduct]: invalid product")

		log.WithError(err).Warn("Invalid product data format")

		return c.JSON(http.StatusBadRequest, entity.ResponseError{
			Error: err.Error(),
		})
	}

	req := entity.ProductManagementRequest{
		ShopID:    shop.ID,
		ProductID: uint32(productID),
	}

	if err := h.usecase.UpdateProduct(&req, &product); err != nil {
		switch err.Error() {
		case "[ShopUsecase.UpdateProduct]: shop not found":
			err = errors.Wrap(err, "[Handler.UpdateProduct]: shop not found")

			log.WithFields(log.Fields{
				"shopID": shop.ID,
			}).WithError(err).Warn("Shop not found")

			return c.JSON(http.StatusNotFound, entity.ResponseError{
				Error: err.Error(),
			})
		case "[ShopUsecase.UpdateProduct]: product not found":
			err = errors.Wrap(err, "[Handler.UpdateProduct]: product not found")

			log.WithFields(log.Fields{
				"productID": productID,
			}).WithError(err).Warn("Product not found")

			return c.JSON(http.StatusNotFound, entity.ResponseError{
				Error: err.Error(),
			})
		case "[ShopUsecase.UpdateProduct]: product does not belong to shop":
			err = errors.Wrap(err, "[Handler.UpdateProduct]: product does not belong to shop")

			log.WithFields(log.Fields{
				"shopID":    shop.ID,
				"productID": productID,
			}).WithError(err).Warn("Product does not belong to shop")

			return c.JSON(http.StatusForbidden, entity.ResponseError{
				Error: err.Error(),
			})
		default:
			err = errors.Wrap(err, "[Handler.UpdateProduct]: internal server error")

			log.WithFields(log.Fields{
				"shopID":    shop.ID,
				"productID": productID,
			}).WithError(err).Error("Internal server error while updating product")

			return c.JSON(http.StatusInternalServerError, entity.ResponseError{
				Error: err.Error(),
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
	shopClaims, ok := c.Get("shop").(*entity.ShopJWT)
	if !ok {
		err := errors.New("[Handler.CreateProduct]: no shop claims found")

		log.Warn("No shop claims found in context")

		return c.JSON(http.StatusUnauthorized, entity.ResponseError{
			Error: err.Error(),
		})
	}

	req := entity.Product{}
	if err := c.Bind(&req); err != nil {
		err = errors.Wrap(err, "[Handler.CreateProduct]: invalid product")

		log.WithError(err).Warn("Invalid product data format")

		return c.JSON(http.StatusBadRequest, entity.ResponseError{
			Error: err.Error(),
		})
	}

	if err := h.usecase.CreateProduct(req, shopClaims.ID); err != nil {
		err = errors.Wrap(err, "[Handler.CreateProduct]: internal server error")

		log.WithFields(log.Fields{
			"shopID": shopClaims.ID,
		}).WithError(err).Error("Internal server error while creating product")

		return c.JSON(http.StatusInternalServerError, entity.ResponseError{
			Error: err.Error(),
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
			err = errors.Wrap(err, "[Handler.GetAllShops]: no shops found")

			log.WithError(err).Warn("No shops found")

			return c.JSON(http.StatusNotFound, entity.ResponseError{
				Error: err.Error(),
			})
		}
		err = errors.Wrap(err, "[Handler.GetAllShops]: internal server error")

		log.WithError(err).Error("Internal server error while getting all shops")

		return c.JSON(http.StatusInternalServerError, entity.ResponseError{
			Error: err.Error(),
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
		err = errors.Wrap(err, "[Handler.CreateShop]: invalid shop")

		log.WithError(err).Warn("Invalid shop data format")

		return c.JSON(http.StatusBadRequest, entity.ResponseError{
			Error: err.Error(),
		})
	}

	if req.Password == "" {
		err := errors.New("[Handler.CreateShop]: password is required")

		log.Warn("Password is required")

		return c.JSON(http.StatusBadRequest, entity.ResponseError{
			Error: err.Error(),
		})
	}

	if err := h.usecase.CreateShop(req); err != nil {
		err = errors.Wrap(err, "[Handler.CreateShop]: internal server error")

		log.WithFields(log.Fields{
			"shopName": req.Name,
		}).WithError(err).Error("Internal server error while creating shop")

		return c.JSON(http.StatusInternalServerError, entity.ResponseError{
			Error: err.Error(),
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
		err = errors.Wrap(err, "[Handler.Login]: invalid shop")

		log.WithError(err).Warn("Invalid login data format")

		return c.JSON(http.StatusBadRequest, entity.ResponseError{
			Error: err.Error(),
		})
	}

	if req.Name == "" || req.Password == "" {
		err := errors.New("[Handler.Login]: name and password are required")

		log.Warn("Name and password are required")

		return c.JSON(http.StatusBadRequest, entity.ResponseError{
			Error: err.Error(),
		})
	}

	token, err := h.usecase.Login(req.Name, req.Password)
	if err != nil {
		if err.Error() == "[ShopUsecase.Login]: shop not found" {
			err = errors.Wrap(err, "[Handler.Login]: shop not found")

			log.WithFields(log.Fields{
				"shopName": req.Name,
			}).WithError(err).Warn("Shop not found during login")

			return c.JSON(http.StatusNotFound, entity.ResponseError{
				Error: err.Error(),
			})
		}
		if err.Error() == "[ShopUsecase.Login]: invalid password" {
			err = errors.Wrap(err, "[Handler.Login]: invalid password")

			log.WithFields(log.Fields{
				"shopName": req.Name,
			}).WithError(err).Warn("Invalid password during login")

			return c.JSON(http.StatusUnauthorized, entity.ResponseError{
				Error: err.Error(),
			})
		}
		err = errors.Wrap(err, "[Handler.Login]: internal server error")

		log.WithFields(log.Fields{
			"shopName": req.Name,
		}).WithError(err).Error("Internal server error during login")

		return c.JSON(http.StatusInternalServerError, entity.ResponseError{
			Error: err.Error(),
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
	shopClaims, ok := c.Get("shop").(*entity.ShopJWT)
	if !ok {
		err := errors.New("[Handler.ReadToken]: no shop claims found")

		log.Warn("No shop claims found in context")

		return c.JSON(http.StatusUnauthorized, entity.ResponseError{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, shopClaims)
}

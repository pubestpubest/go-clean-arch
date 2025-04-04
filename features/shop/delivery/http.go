package delivery

import (
	"net/http"
	"order-management/domain"
	"order-management/entity"
	"strconv"

	"order-management/middleware"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	usecase domain.ShopUsecase
}

func NewHandler(e *echo.Group, u domain.ShopUsecase) *Handler {
	h := Handler{usecase: u}

	// Public group - no authentication required
	publicGroup := e.Group("")
	publicGroup.GET("/", h.GetAllShops)                          // Anyone can view shops
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
	claims, ok := c.Get("shop").(*entity.ShopResponse)
	if !ok {
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}
	shop, err := h.usecase.GetShopByName(claims.Name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, shop)
}

func (h *Handler) GetProductsByShopID(c echo.Context) error {
	shopID, err := strconv.ParseUint(c.Param("shop_id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	products, err := h.usecase.GetProductsByShopID(uint32(shopID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, products)
}

func (h *Handler) DeleteProduct(c echo.Context) error {
	productID, err := strconv.ParseUint(c.Param("product_id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	shop, ok := c.Get("shop").(*entity.ShopResponse)
	if !ok {
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}
	req := entity.ProductManagementRequest{
		ShopResponse: *shop,
		ProductID:    uint32(productID),
	}
	if err := h.usecase.DeleteProduct(&req); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusOK)
}

func (h *Handler) UpdateProduct(c echo.Context) error {
	//Param
	productID, err := strconv.ParseUint(c.Param("product_id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	//Get shop from JWT
	shop, ok := c.Get("shop").(*entity.ShopResponse)
	if !ok {
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}
	//Bind
	product := entity.Product{}
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	req := entity.ProductManagementRequest{
		ShopResponse: *shop,
		ProductID:    uint32(productID),
	}
	//Update product
	if err := h.usecase.UpdateProduct(&req, &product); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusOK)
}

func (h *Handler) Logout(c echo.Context) error {
	cookie := &http.Cookie{
		Name:   "token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	c.SetCookie(cookie)
	return c.NoContent(http.StatusOK)
}

func (h *Handler) CreateProduct(c echo.Context) error {
	claims, ok := c.Get("shop").(*entity.ShopResponse)
	if !ok {
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}
	req := entity.Product{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := h.usecase.CreateProduct(req, claims.ID); err != nil {
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
	if req.Password == "" {
		return c.JSON(http.StatusBadRequest, "password is required")
	}
	if err := h.usecase.CreateShop(req); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) Login(c echo.Context) error {
	req := entity.Shop{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if req.Name == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, "name and password are required")
	}

	token, err := h.usecase.Login(req.Name, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}
	c.Response().Header().Set("Authorization", "Bearer "+token)
	return c.NoContent(http.StatusOK)
}

func (h *Handler) ReadToken(c echo.Context) error {
	claims, ok := c.Get("shop").(*entity.ShopResponse)
	if !ok {
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}
	return c.JSON(http.StatusOK, claims)
}

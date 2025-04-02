package delivery

import (
	"net/http"
	"order-management/domain"
	"order-management/entity"
	"order-management/response"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	usecase domain.ShopUsecase
}

func NewHandler(e *echo.Group, u domain.ShopUsecase) *Handler {
	h := Handler{usecase: u}

	// Need JWT
	e.POST("/products", h.CreateProduct)               // Only shop owner can create products
	e.PUT("/products/:product_id", h.UpdateProduct)    // Only shop owner can update their products
	e.DELETE("/products/:product_id", h.DeleteProduct) // Only shop owner can delete their products
	e.GET("/shops/me", h.ReadToken)                    // Get current shop profile from JWT
	e.POST("/shops/logout", h.Logout)                  // Logout requires JWT
	e.GET("/shops/profile", h.GetShopProfile)          // Get detailed profile requires JWT

	// Public endpoints
	e.GET("/:shop_id/products", h.GetProducts)               // Anyone can view products
	e.GET("/shops", h.GetAllShops)                           // Anyone can view shops
	e.POST("/shops/register", h.CreateShop)                  // Public registration
	e.POST("/shops/login", h.Login)                          // Public login
	e.GET("/shops/:shop_id/products", h.GetProductsByShopID) // Anyone can view products
	return &h
}

func (h *Handler) GetShopProfile(c echo.Context) error {
	claims, err := readToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
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
	claims, err := readToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}
	if !h.usecase.BelongsToShop(uint32(productID), claims) {
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}
	if err := h.usecase.DeleteProduct(uint32(productID)); err != nil {
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
	//Read token
	claims, err := readToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}
	//Bind
	req := entity.Product{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	//Check if the product belongs to the shop
	if !h.usecase.BelongsToShop(uint32(productID), claims) {
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}
	//Update product
	if err := h.usecase.UpdateProduct(uint32(productID), &req); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusOK)
}

func (h *Handler) GetProducts(c echo.Context) error {
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
	claims, err := readToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
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

	shop, err := h.usecase.Login(req.Name, req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if err := bcrypt.CompareHashAndPassword([]byte(shop.Password), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, "invalid password")
	}

	products, err := h.usecase.GetProductsByShopID(shop.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	var data = map[string]interface{}{
		"id":          shop.ID,
		"name":        shop.Name,
		"description": shop.Description,
		"products":    products,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":          shop.ID,
		"name":        shop.Name,
		"description": shop.Description,
	})
	t, err := token.SignedString([]byte(viper.GetString("jwt.secret")))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	cookie := &http.Cookie{
		Name:  "token",
		Value: t,
		Path:  "/",
	}
	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, data)
}

func (h *Handler) ReadToken(c echo.Context) error {
	claims, err := readToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}
	return c.JSON(http.StatusOK, claims)
}

func readToken(c echo.Context) (data *response.Shop, err error) {
	cookie, err := c.Cookie("token")
	if err != nil {
		return nil, err
	}
	token, err := jwt.ParseWithClaims(cookie.Value, &response.Shop{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("jwt.secret")), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*response.Shop)
	if !ok || !token.Valid {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	data = &response.Shop{
		ID:          claims.ID,
		Name:        claims.Name,
		Description: claims.Description,
	}
	return data, nil
}

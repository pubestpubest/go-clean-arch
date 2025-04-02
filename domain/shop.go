package domain

import (
	"order-management/entity"
	"order-management/response"
)

type ShopUsecase interface {
	CreateProduct(product entity.Product, shopID uint32) error
	CreateShop(shop entity.Shop) error
	GetAllShopsWithProducts() ([]response.ShopWithProducts, error)
	GetAllShops() ([]response.Shop, error)
	GetShopByName(name string) (entity.Shop, error)
	Login(name string, password string) (entity.Shop, error)
	GetProductsByShopID(id uint32) ([]response.Product, error)
	BelongsToShop(productID uint32, claims *response.Shop) bool
	UpdateProduct(productID uint32, newProduct *entity.Product) error
	DeleteProduct(productID uint32) error
}

type ShopRepository interface {
	CreateProduct(product entity.Product, shopID uint32) error
	CreateShop(shop entity.Shop) error
	GetAllShops() ([]entity.Shop, error)
	GetProductsByShopID(shopID uint32) ([]entity.Product, error)
	GetShopByName(name string) (entity.Shop, error)
	UpdateProduct(productID uint32, newProduct *entity.Product) error
	GetProductByID(productID uint32) (entity.Product, error)
	DeleteProduct(productID uint32) error
}

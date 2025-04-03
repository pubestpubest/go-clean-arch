package domain

import (
	"order-management/entity"
)

type ShopUsecase interface {
	CreateProduct(product entity.Product, shopID uint32) error
	CreateShop(shop entity.Shop) error
	GetAllShopsWithProducts() ([]entity.ShopWithProducts, error)
	GetAllShops() ([]entity.Shop, error)
	GetShopByName(name string) (entity.ShopWithProducts, error)
	Login(name string, password string) (entity.Shop, error)
	GetProductsByShopID(id uint32) ([]entity.Product, error)
	BelongsToShop(productID uint32, claims *entity.ShopResponse) bool
	UpdateProduct(productID uint32, newProduct *entity.Product) error
	DeleteProduct(productID uint32) error
}

type ShopRepository interface {
	CreateProduct(product entity.Product, shopID uint32) error
	CreateShop(shop entity.Shop) error
	GetAllShops() ([]entity.Shop, error)
	GetProductsByShopID(shopID uint32) ([]entity.Product, error)
	GetShopByName(name string) (entity.ShopResponse, error)
	UpdateProduct(productID uint32, newProduct *entity.Product) error
	GetProductByID(productID uint32) (entity.Product, error)
	DeleteProduct(productID uint32) error
	GetShopByNameWithPassword(name string) (entity.Shop, error)
}

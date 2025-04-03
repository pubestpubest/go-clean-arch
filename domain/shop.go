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
	Login(name string, password string) (string, error)
	GetProductsByShopID(id uint32) ([]entity.Product, error)
	UpdateProduct(req *entity.ProductManagementRequest, product *entity.Product) error
	DeleteProduct(req *entity.ProductManagementRequest) error
}

type ShopRepository interface {
	CreateProduct(product entity.Product, shopID uint32) error
	CreateShop(shop entity.Shop) error
	GetAllShops() ([]entity.Shop, error)
	GetProductsByShopID(shopID uint32) ([]entity.Product, error)
	GetShopByName(name string) (entity.ShopResponse, error)
	UpdateProduct(req *entity.ProductManagementRequest, product *entity.Product) error
	GetProductByID(productID uint32) (entity.Product, error)
	DeleteProduct(req *entity.ProductManagementRequest) error
	GetShopByNameWithPassword(name string) (entity.Shop, error)
}

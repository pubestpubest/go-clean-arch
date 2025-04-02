package domain

import (
	"order-management/entity"
	"order-management/response"
)

type ShopUsecase interface {
	CreateProduct(product entity.Product, shopID uint32) error
	CreateShop(shop entity.Shop) error
	GetAllShopsWithProducts() ([]response.Shop, error)
	GetShopByName(name string) (entity.Shop, error)
	Login(name string, password string) (entity.Shop, error)
	GetProductsByShopID(id uint32) ([]response.Product, error)
}

type ShopRepository interface {
	CreateProduct(product entity.Product, shopID uint32) error
	CreateShop(shop entity.Shop) error
	GetAllShops() ([]entity.Shop, error)
	GetProductsByShopID(shopID uint32) ([]entity.Product, error)
	GetShopByName(name string) (entity.Shop, error)
}

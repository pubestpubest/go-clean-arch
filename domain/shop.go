package domain

import "order-management/entity"

type ShopUsecase interface {
	CreateProduct(product entity.Product, shopID uint32) error
	CreateShop(shop entity.Shop) error
	GetAllShops() ([]entity.Shop, error)
}

type ShopRepository interface {
	CreateProduct(product entity.Product, shopID uint32) error
	CreateShop(shop entity.Shop) error
	GetAllShops() ([]entity.Shop, error)
}

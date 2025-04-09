package domain

import (
	"order-management/entity"
)

type ProductUsecase interface {
	GetAllProducts() ([]entity.ProductWithOutShop, error)
	GetProductPrice(productID uint32) (float64, error)
}

type ProductRepository interface {
	CreateProduct(product entity.Product, shopID uint32) error
	GetProductsByShopID(shopID uint32) ([]entity.ProductWithOutShop, error)
	UpdateProduct(req *entity.ProductManagementRequest, product *entity.Product) error
	GetProductByID(productID uint32) (entity.Product, error)
	DeleteProduct(req *entity.ProductManagementRequest) error
	GetAllProducts() ([]entity.ProductWithOutShop, error)
	GetProductPrice(productID uint32) (float64, error)
}

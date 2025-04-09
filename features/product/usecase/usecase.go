package usecase

import (
	"order-management/domain"
	"order-management/entity"
)

type productUsecase struct {
	productRepo domain.ProductRepository
}

func NewProductUsecase(productRepo domain.ProductRepository) domain.ProductUsecase {
	return &productUsecase{
		productRepo: productRepo,
	}
}

func (u *productUsecase) GetAllProducts() ([]entity.ProductWithOutShop, error) {
	return u.productRepo.GetAllProducts()
}

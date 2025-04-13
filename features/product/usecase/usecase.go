package usecase

import (
	"order-management/domain"
	"order-management/entity"

	"github.com/pkg/errors"
)

type productUsecase struct {
	productRepo domain.ProductRepository
}

func NewProductUsecase(productRepo domain.ProductRepository) domain.ProductUsecase {
	return &productUsecase{
		productRepo: productRepo,
	}
}

func (u *productUsecase) GetProductByID(productID uint32) (entity.Product, error) {
	product, err := u.productRepo.GetProductByID(productID)
	if err != nil {
		err = errors.Wrap(err, "[ProductUsecase.GetProductByID]: failed to get product by id")
		return entity.Product{}, err
	}
	return product, nil
}

func (u *productUsecase) GetAllProducts() ([]entity.ProductWithOutShop, error) {
	products, err := u.productRepo.GetAllProducts()
	if err != nil {
		err = errors.Wrap(err, "[ProductUsecase.GetAllProducts]: failed to get all products")
		return nil, err
	}
	return products, nil
}

func (u *productUsecase) GetProductPrice(productID uint32) (float64, error) {
	price, err := u.productRepo.GetProductPrice(productID)
	if err != nil {
		err = errors.Wrap(err, "[ProductUsecase.GetProductPrice]: failed to get product price")
		return 0, err
	}
	return price, nil
}

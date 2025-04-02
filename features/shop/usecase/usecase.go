package usecase

import (
	"order-management/domain"
	"order-management/entity"
)

type shopUsecase struct {
	repo domain.ShopRepository
}

func NewShopUsecase(repo domain.ShopRepository) domain.ShopUsecase {
	return &shopUsecase{repo: repo}
}

func (u *shopUsecase) CreateProduct(product entity.Product, shopID uint32) error {
	return u.repo.CreateProduct(product, shopID)
}

func (u *shopUsecase) CreateShop(shop entity.Shop) error {
	return u.repo.CreateShop(shop)
}

func (u *shopUsecase) GetAllShops() ([]entity.Shop, error) {
	return u.repo.GetAllShops()
}

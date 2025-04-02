package repository

import (
	"order-management/domain"
	"order-management/entity"

	"gorm.io/gorm"
)

type shopRepository struct {
	db *gorm.DB
}

func NewShopRepository(db *gorm.DB) domain.ShopRepository {
	return &shopRepository{db: db}
}

func (r *shopRepository) CreateProduct(product entity.Product, shopID uint32) error {
	product.ShopID = shopID
	return r.db.Create(&product).Error
}

func (r *shopRepository) CreateShop(shop entity.Shop) error {
	return r.db.Create(&shop).Error
}

func (r *shopRepository) GetAllShops() ([]entity.Shop, error) {
	var shops []entity.Shop
	if err := r.db.Preload("Products").Find(&shops, 1).Error; err != nil {
		return nil, err
	}
	return shops, nil
}

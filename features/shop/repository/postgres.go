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

func (r *shopRepository) GetAllShops() (shops []entity.Shop, err error) {
	if err := r.db.Find(&shops).Error; err != nil {
		return nil, err
	}
	return shops, nil
}

func (r *shopRepository) GetProductsByShopID(shopID uint32) (products []entity.Product, err error) {
	if err := r.db.Where("shop_id = ?", shopID).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

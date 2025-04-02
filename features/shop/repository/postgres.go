package repository

import (
	"order-management/domain"
	"order-management/entity"
	"order-management/response"

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

func (r *shopRepository) GetShopByName(name string) (shop response.Shop, err error) {
	var entityShop entity.Shop
	if err := r.db.Where("name = ?", name).First(&entityShop).Error; err != nil {
		return response.Shop{}, err
	}
	shop = response.Shop{
		ID:          entityShop.ID,
		Name:        entityShop.Name,
		Description: entityShop.Description,
	}
	return shop, nil
}

func (r *shopRepository) GetShopByNameWithPassword(name string) (shop entity.Shop, err error) {
	if err := r.db.Where("name = ?", name).First(&shop).Error; err != nil {
		return entity.Shop{}, err
	}
	return shop, nil
}

func (r *shopRepository) UpdateProduct(productID uint32, newProduct *entity.Product) error {
	return r.db.Model(&entity.Product{}).Where("id = ?", productID).Updates(newProduct).Error
}

func (r *shopRepository) GetProductByID(productID uint32) (product entity.Product, err error) {
	if err := r.db.First(&product, productID).Error; err != nil {
		return entity.Product{}, err
	}
	return product, nil
}

func (r *shopRepository) DeleteProduct(productID uint32) error {
	return r.db.Where("id = ?", productID).Delete(&entity.Product{}).Error
}

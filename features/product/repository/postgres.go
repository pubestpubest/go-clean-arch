package repository

import (
	"order-management/domain"
	"order-management/entity"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetProductPrice(productID uint32) (float64, error) {
	var product entity.Product
	if err := r.db.Where("id = ?", productID).First(&product).Error; err != nil {
		return 0, errors.Wrap(err, "[ProductRepository.GetProductPrice]: failed to get product price")
	}
	return float64(product.Price), nil
}

func (r *productRepository) CreateProduct(product entity.Product, shopID uint32) error {
	product.ShopID = shopID
	if err := r.db.Create(&product).Error; err != nil {
		err = errors.Wrap(err, "[ProductRepository.CreateProduct]: failed to create product")
		return err
	}
	return nil
}

func (r *productRepository) GetProductsByShopID(shopID uint32) (products []entity.ProductWithOutShop, err error) {
	if err := r.db.Model(&entity.Product{}).Where("shop_id = ?", shopID).Find(&products).Error; err != nil {
		err = errors.Wrap(err, "[ProductRepository.GetProductsByShopID]: failed to get products by shop id")
		return nil, err
	}
	return products, nil
}

func (r *productRepository) UpdateProduct(req *entity.ProductManagementRequest, product *entity.Product) error {
	if err := r.db.Model(&entity.Product{}).Where("id = ? AND shop_id = ?", req.ProductID, req.ShopID).Updates(product).Error; err != nil {
		err = errors.Wrap(err, "[ProductRepository.UpdateProduct]: failed to update product")
		return err
	}
	return nil
}

func (r *productRepository) GetProductByID(productID uint32) (product entity.Product, err error) {
	if err := r.db.Preload("Shop").First(&product, productID).Error; err != nil {
		err = errors.Wrap(err, "[ProductRepository.GetProductByID]: failed to get product by id")
		return entity.Product{}, err
	}
	product.Shop.Password = ""
	return product, nil
}

func (r *productRepository) DeleteProduct(req *entity.ProductManagementRequest) error {
	if err := r.db.Where("id = ? AND shop_id = ?", req.ProductID, req.ShopID).Delete(&entity.Product{}).Error; err != nil {
		err = errors.Wrap(err, "[ProductRepository.DeleteProduct]: failed to delete product")
		return err
	}
	return nil
}

func (r *productRepository) GetAllProducts() (products []entity.ProductWithOutShop, err error) {
	if err := r.db.Model(&entity.Product{}).Find(&products).Error; err != nil {
		err = errors.Wrap(err, "[ProductRepository.GetAllProducts]: failed to get all products")
		return nil, err
	}
	return products, nil
}

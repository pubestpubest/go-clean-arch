package repository

import (
	"order-management/domain"
	"order-management/entity"

	"github.com/pkg/errors"
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
	if err := r.db.Create(&product).Error; err != nil {
		return errors.Wrap(err, "[ShopRepository.CreateProduct]: failed to create product")
	}
	return nil
}

func (r *shopRepository) CreateShop(shop entity.Shop) error {
	if err := r.db.Create(&shop).Error; err != nil {
		return errors.Wrap(err, "[ShopRepository.CreateShop]: failed to create shop")
	}
	return nil
}

func (r *shopRepository) GetAllShops() (shops []entity.Shop, err error) {
	if err := r.db.Find(&shops).Error; err != nil {
		return nil, errors.Wrap(err, "[ShopRepository.GetAllShops]: failed to get all shops")
	}
	return shops, nil
}

func (r *shopRepository) GetProductsByShopID(shopID uint32) (products []entity.Product, err error) {
	if err := r.db.Where("shop_id = ?", shopID).Find(&products).Error; err != nil {
		return nil, errors.Wrap(err, "[ShopRepository.GetProductsByShopID]: failed to get products by shop id")
	}
	return products, nil
}

func (r *shopRepository) GetShopByName(name string) (shop entity.ShopWithOutPassword, err error) {
	var entityShop entity.Shop
	if err := r.db.Where("name = ?", name).First(&entityShop).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.ShopWithOutPassword{}, errors.New("[ShopRepository.GetShopByName]: shop not found")
		}
		return entity.ShopWithOutPassword{}, errors.Wrap(err, "[ShopRepository.GetShopByName]: failed to get shop by name")
	}
	shop = entity.ShopWithOutPassword{
		ID:          entityShop.ID,
		Name:        entityShop.Name,
		Description: entityShop.Description,
	}
	return shop, nil
}

func (r *shopRepository) GetShopByNameWithPassword(name string) (shop entity.Shop, err error) {
	if err := r.db.Where("name = ?", name).First(&shop).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Shop{}, errors.New("[ShopRepository.GetShopByNameWithPassword]: shop not found")
		}
		return entity.Shop{}, errors.Wrap(err, "[ShopRepository.GetShopByNameWithPassword]: failed to get shop by name with password")
	}
	return shop, nil
}

func (r *shopRepository) UpdateProduct(req *entity.ProductManagementRequest, newProduct *entity.Product) error {
	if err := r.db.Model(&entity.Product{}).Where("id = ? AND shop_id = ?", req.ProductID, req.ShopID).Updates(newProduct).Error; err != nil {
		return errors.Wrap(err, "[ShopRepository.UpdateProduct]: failed to update product")
	}
	return nil
}

func (r *shopRepository) GetProductByID(productID uint32) (product entity.Product, err error) {
	if err := r.db.First(&product, productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Product{}, errors.New("[ShopRepository.GetProductByID]: product not found")
		}
		return entity.Product{}, errors.Wrap(err, "[ShopRepository.GetProductByID]: failed to get product by id")
	}
	return product, nil
}

func (r *shopRepository) DeleteProduct(req *entity.ProductManagementRequest) error {
	if err := r.db.Where("id = ? AND shop_id = ?", req.ProductID, req.ShopID).Delete(&entity.Product{}).Error; err != nil {
		return errors.Wrap(err, "[ShopRepository.DeleteProduct]: failed to delete product")
	}
	return nil
}

func (r *shopRepository) ShopExists(id uint32) (bool, error) {
	var exists bool
	err := r.db.Model(&entity.Shop{}).
		Select("count(*) > 0").
		Where("id = ?", id).
		Find(&exists).
		Error
	if err != nil {
		return false, errors.Wrap(err, "[ShopRepository.ShopExists]: failed to check shop existence")
	}
	return exists, nil
}

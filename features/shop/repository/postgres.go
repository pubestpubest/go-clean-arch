package repository

import (
	"order-management/domain"
	"order-management/entity"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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
		err = errors.Wrap(err, "[ShopRepository.CreateProduct]: failed to create product")
		log.WithFields(log.Fields{
			"product": product,
			"shopID":  shopID,
		}).WithError(err).Error("failed to create product")
		return err
	}
	return nil
}

func (r *shopRepository) CreateShop(shop entity.Shop) error {
	if err := r.db.Create(&shop).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			err = errors.New("[ShopRepository.CreateShop]: shop already exists")
			log.WithFields(log.Fields{
				"shop": shop,
			}).WithError(err).Warn("shop already exists")
			return err
		}
		err = errors.Wrap(err, "[ShopRepository.CreateShop]: failed to create shop")
		log.WithFields(log.Fields{
			"shop": shop,
		}).WithError(err).Error("failed to create shop")
		return err
	}
	return nil
}

func (r *shopRepository) GetAllShops() (shops []entity.Shop, err error) {
	if err := r.db.Find(&shops).Error; err != nil {
		err = errors.Wrap(err, "[ShopRepository.GetAllShops]: failed to get all shops")
		log.WithError(err).Error("failed to get all shops")
		return nil, err
	}
	return shops, nil
}

func (r *shopRepository) GetProductsByShopID(shopID uint32) (products []entity.Product, err error) {
	if err := r.db.Where("shop_id = ?", shopID).Find(&products).Error; err != nil {
		err = errors.Wrap(err, "[ShopRepository.GetProductsByShopID]: failed to get products by shop id")
		log.WithFields(log.Fields{
			"shopID": shopID,
		}).WithError(err).Error("failed to get products by shop id")
		return nil, err
	}
	return products, nil
}

func (r *shopRepository) GetShopByName(name string) (shop entity.ShopWithOutPassword, err error) {
	if err := r.db.Select("id", "name", "description").Where("name = ?", name).First(&shop).Error; err != nil {
		err = errors.Wrap(err, "[ShopRepository.GetShopByName]: failed to get shop by name")
		log.WithFields(log.Fields{
			"name": name,
		}).WithError(err).Error("failed to get shop by name")
		return entity.ShopWithOutPassword{}, err
	}
	return shop, nil
}

func (r *shopRepository) GetShopByNameWithPassword(name string) (shop entity.Shop, err error) {
	if err := r.db.Where("name = ?", name).First(&shop).Error; err != nil {
		err = errors.Wrap(err, "[ShopRepository.GetShopByNameWithPassword]: failed to get shop by name with password")
		log.WithFields(log.Fields{
			"name": name,
		}).WithError(err).Error("failed to get shop by name with password")
		return entity.Shop{}, err
	}
	return shop, nil
}

func (r *shopRepository) UpdateProduct(req *entity.ProductManagementRequest, newProduct *entity.Product) error {
	if err := r.db.Model(&entity.Product{}).Where("id = ? AND shop_id = ?", req.ProductID, req.ShopID).Updates(newProduct).Error; err != nil {
		err = errors.Wrap(err, "[ShopRepository.UpdateProduct]: failed to update product")
		log.WithFields(log.Fields{
			"productID": req.ProductID,
			"shopID":    req.ShopID,
		}).WithError(err).Error("failed to update product")
		return err
	}
	return nil
}

func (r *shopRepository) GetProductByID(productID uint32) (product entity.Product, err error) {
	if err := r.db.First(&product, productID).Error; err != nil {
		err = errors.Wrap(err, "[ShopRepository.GetProductByID]: failed to get product by id")
		log.WithFields(log.Fields{
			"productID": productID,
		}).WithError(err).Error("failed to get product by id")
		return entity.Product{}, err
	}
	return product, nil
}

func (r *shopRepository) DeleteProduct(req *entity.ProductManagementRequest) error {
	if err := r.db.Where("id = ? AND shop_id = ?", req.ProductID, req.ShopID).Delete(&entity.Product{}).Error; err != nil {
		err = errors.Wrap(err, "[ShopRepository.DeleteProduct]: failed to delete product")
		log.WithFields(log.Fields{
			"productID": req.ProductID,
			"shopID":    req.ShopID,
		}).WithError(err).Error("failed to delete product")
		return err
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
		err = errors.Wrap(err, "[ShopRepository.ShopExists]: failed to check shop existence")
		log.WithFields(log.Fields{
			"id": id,
		}).WithError(err).Error("failed to check shop existence")
		return false, err
	}
	return exists, nil
}

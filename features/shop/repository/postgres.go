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

// Trace -> Enter, Exit for dev
// Debug -> Detail for dev
// Info -> Status for system cron job
// Warn -> Something wrong but handleable [User's fault]
// Error -> Something wrong didn't handle [Developer's fault]
// Fatal -> Something serious happened, system can't handle it
// Focus to log on the failed case
// Happy case is not that important

// ✅
func (r *shopRepository) CreateProduct(product entity.Product, shopID uint32) error {

	log.Trace("Entering function CreateProduct()")
	defer log.Trace("Exiting function CreateProduct()")

	product.ShopID = shopID

	log.WithFields(log.Fields{
		"product": product,
		"shopID":  shopID,
	}).Debug("Creating product")

	if err := r.db.Create(&product).Error; err != nil {
		err = errors.Wrap(err, "[ShopRepository.CreateProduct]: failed to create product")
		return err
	}

	return nil
}

// ✅
func (r *shopRepository) CreateShop(shop entity.Shop) error {

	log.Trace("Entering function CreateShop()")
	defer log.Trace("Exiting function CreateShop()")

	log.WithFields(log.Fields{
		"shop": shop,
	}).Debug("Creating shop")

	if err := r.db.Create(&shop).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			err = errors.New("[ShopRepository.CreateShop]: shop already exists")

			log.WithFields(log.Fields{
				"shop": shop,
			}).WithError(err).Warn("shop already exists")

			return err
		}
		err = errors.Wrap(err, "[ShopRepository.CreateShop]: failed to create shop")
		return err
	}

	return nil
}

// ✅
func (r *shopRepository) GetAllShops() (shops []entity.Shop, err error) {

	log.Trace("Entering function GetAllShops()")
	defer log.Trace("Exiting function GetAllShops()")

	if err := r.db.Find(&shops).Error; err != nil {
		err = errors.Wrap(err, "[ShopRepository.GetAllShops]: failed to get all shops")
		return nil, err
	}

	return shops, nil
}

// ✅
func (r *shopRepository) GetProductsByShopID(shopID uint32) (products []entity.Product, err error) {

	log.Trace("Entering function GetProductsByShopID()")
	defer log.Trace("Exiting function GetProductsByShopID()")

	log.WithFields(log.Fields{
		"shopID": shopID,
	}).Debug("Getting products by shop id")

	if err := r.db.Where("shop_id = ?", shopID).Find(&products).Error; err != nil {
		err = errors.Wrap(err, "[ShopRepository.GetProductsByShopID]: failed to get products by shop id")
		return nil, err
	}

	return products, nil
}

// ✅
func (r *shopRepository) GetShopByName(name string) (shop entity.ShopWithOutPassword, err error) {
	log.Trace("Entering function GetShopByName()")
	defer log.Trace("Exiting function GetShopByName()")

	log.WithFields(log.Fields{
		"name": name,
	}).Debug("Getting shop by name")

	if err := r.db.Model(&entity.Shop{}).Select("id", "name", "description").Where("name = ?", name).First(&shop).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("[ShopRepository.GetShopByName]: shop not found")
			return entity.ShopWithOutPassword{}, err
		}
		err = errors.Wrap(err, "[ShopRepository.GetShopByName]: failed to get shop by name")
		return entity.ShopWithOutPassword{}, err
	}

	return shop, nil
}

// ✅
func (r *shopRepository) GetShopByNameWithPassword(name string) (shop entity.Shop, err error) {
	log.Trace("Entering function GetShopByNameWithPassword()")
	defer log.Trace("Exiting function GetShopByNameWithPassword()")

	log.WithFields(log.Fields{
		"name": name,
	}).Debug("Getting shop by name with password")

	if err := r.db.Where("name = ?", name).First(&shop).Error; err != nil {
		err = errors.Wrap(err, "[ShopRepository.GetShopByNameWithPassword]: failed to get shop by name with password")
		return entity.Shop{}, err
	}

	return shop, nil
}

// ✅
func (r *shopRepository) UpdateProduct(req *entity.ProductManagementRequest, newProduct *entity.Product) error {
	log.Trace("Entering function UpdateProduct()")
	defer log.Trace("Exiting function UpdateProduct()")

	log.WithFields(log.Fields{
		"productID": req.ProductID,
		"shopID":    req.ShopID,
	}).Debug("Updating product")

	if err := r.db.Model(&entity.Product{}).Where("id = ? AND shop_id = ?", req.ProductID, req.ShopID).Updates(newProduct).Error; err != nil {
		err = errors.Wrap(err, "[ShopRepository.UpdateProduct]: failed to update product")
		return err
	}

	return nil
}

// ✅
func (r *shopRepository) GetProductByID(productID uint32) (product entity.Product, err error) {
	log.Trace("Entering function GetProductByID()")
	defer log.Trace("Exiting function GetProductByID()")

	log.WithFields(log.Fields{
		"productID": productID,
	}).Debug("Getting product by id")

	if err := r.db.First(&product, productID).Error; err != nil {
		err = errors.Wrap(err, "[ShopRepository.GetProductByID]: failed to get product by id")
		return entity.Product{}, err
	}

	return product, nil
}

// ✅
func (r *shopRepository) DeleteProduct(req *entity.ProductManagementRequest) error {
	log.Trace("Entering function DeleteProduct()")
	defer log.Trace("Exiting function DeleteProduct()")

	log.WithFields(log.Fields{
		"productID": req.ProductID,
		"shopID":    req.ShopID,
	}).Debug("Deleting product")

	if err := r.db.Where("id = ? AND shop_id = ?", req.ProductID, req.ShopID).Delete(&entity.Product{}).Error; err != nil {
		err = errors.Wrap(err, "[ShopRepository.DeleteProduct]: failed to delete product")
		return err
	}

	return nil
}

// ✅
func (r *shopRepository) ShopExists(id uint32) (bool, error) {
	log.Trace("Entering function ShopExists()")
	defer log.Trace("Exiting function ShopExists()")

	log.WithFields(log.Fields{
		"id": id,
	}).Debug("Checking shop existence")

	var exists bool
	err := r.db.Model(&entity.Shop{}).
		Select("count(*) > 0").
		Where("id = ?", id).
		Find(&exists).
		Error
	if err != nil {
		err = errors.Wrap(err, "[ShopRepository.ShopExists]: failed to check shop existence")
		return false, err
	}

	return exists, nil
}

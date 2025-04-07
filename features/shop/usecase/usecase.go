package usecase

import (
	"order-management/domain"
	"order-management/entity"
	"order-management/utils"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type shopUsecase struct {
	repo domain.ShopRepository
}

func NewShopUsecase(repo domain.ShopRepository) domain.ShopUsecase {
	return &shopUsecase{repo: repo}
}

func (u *shopUsecase) CreateProduct(product entity.Product, shopID uint32) error {
	log.Trace("Entering function CreateProduct()")
	defer log.Trace("Exiting function CreateProduct()")

	log.WithFields(log.Fields{
		"product": product,
		"shopID":  shopID,
	}).Debug("Creating product")

	if err := u.repo.CreateProduct(product, shopID); err != nil {
		err = errors.Wrap(err, "[ShopUsecase.CreateProduct]: failed to create product")
		return err
	}

	return nil
}

func (u *shopUsecase) CreateShop(shop entity.Shop) error {
	log.Trace("Entering function CreateShop()")
	defer log.Trace("Exiting function CreateShop()")

	log.WithFields(log.Fields{
		"shop": shop,
	}).Debug("Creating shop")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(shop.Password), bcrypt.DefaultCost)
	if err != nil {
		err = errors.Wrap(err, "[ShopUsecase.CreateShop]: failed to hash password")
		return err
	}

	shop.Password = string(hashedPassword)

	if err := u.repo.CreateShop(shop); err != nil {
		if err.Error() == "[ShopRepository.CreateShop]: shop already exists" {
			return err
		}
		err = errors.Wrap(err, "[ShopUsecase.CreateShop]: failed to create shop")
		return err
	}
	return nil
}

// Too big O(n^2)
func (u *shopUsecase) GetAllShopsWithProducts() ([]entity.ShopWithProducts, error) {
	log.Trace("Entering function GetAllShopsWithProducts()")
	defer log.Trace("Exiting function GetAllShopsWithProducts()")

	shops, err := u.repo.GetAllShops()
	if err != nil {
		err = errors.Wrap(err, "[ShopUsecase.GetAllShopsWithProducts]: failed to get all shops")
		return nil, err
	}

	if len(shops) == 0 {
		err = errors.Wrap(gorm.ErrRecordNotFound, "[ShopUsecase.GetAllShopsWithProducts]: no shops found")
		return nil, err
	}

	shopsResponse := []entity.ShopWithProducts{}
	for _, shop := range shops {
		products, err := u.repo.GetProductsByShopID(shop.ID)
		if err != nil {
			err = errors.Wrap(err, "[ShopUsecase.GetAllShopsWithProducts]: failed to get products by shop id")
			return nil, err
		}
		productsResponse := []entity.ProductResponse{}
		for _, product := range products {
			productsResponse = append(productsResponse, entity.ProductResponse{
				ID:          product.ID,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
			})
		}
		shopResponse := entity.ShopWithProducts{
			ID:          shop.ID,
			Name:        shop.Name,
			Description: shop.Description,
			Products:    productsResponse,
		}
		shopsResponse = append(shopsResponse, shopResponse)
	}

	return shopsResponse, nil
}

func (u *shopUsecase) GetAllShops() ([]entity.Shop, error) {
	log.Trace("Entering function GetAllShops()")
	defer log.Trace("Exiting function GetAllShops()")

	log.Debug("Getting all shops")

	shops, err := u.repo.GetAllShops()
	if err != nil {
		err = errors.Wrap(err, "[ShopUsecase.GetAllShops]: failed to get all shops")
		return nil, err
	}

	shopsResponse := []entity.Shop{}
	for _, shop := range shops {
		shopsResponse = append(shopsResponse, entity.Shop{
			ID:          shop.ID,
			Name:        shop.Name,
			Description: shop.Description,
		})
	}

	return shopsResponse, nil
}

func (u *shopUsecase) GetShopByName(name string) (entity.ShopWithProducts, error) {
	log.Trace("Entering function GetShopByName()")
	defer log.Trace("Exiting function GetShopByName()")

	log.Debug("Getting shop by name")

	shop, err := u.repo.GetShopByName(name)
	if err != nil {
		if err.Error() == "[ShopRepository.GetShopByName]: shop not found" {
			err = errors.New("[ShopUsecase.GetShopByName]: shop not found")
			return entity.ShopWithProducts{}, err
		}

		err = errors.Wrap(err, "[ShopUsecase.GetShopByName]: failed to get shop by name")
		return entity.ShopWithProducts{}, err
	}

	products, err := u.repo.GetProductsByShopID(shop.ID)
	if err != nil {
		err = errors.Wrap(err, "[ShopUsecase.GetShopByName]: failed to get products by shop id")
		return entity.ShopWithProducts{}, err
	}

	productsResponse := []entity.ProductResponse{}
	for _, product := range products {
		productsResponse = append(productsResponse, entity.ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}

	shopResponse := entity.ShopWithProducts{
		ID:          shop.ID,
		Name:        shop.Name,
		Description: shop.Description,
		Products:    productsResponse,
	}

	return shopResponse, nil
}

func (u *shopUsecase) Login(name string, password string) (string, error) {
	log.Trace("Entering function Login()")
	defer log.Trace("Exiting function Login()")

	log.WithFields(log.Fields{
		"name": name,
	}).Debug("Logging in")

	credentials, err := u.repo.GetShopByNameWithPassword(name)
	if err != nil {
		if err.Error() == "[ShopRepository.GetShopByNameWithPassword]: shop not found" {
			err = errors.New("[ShopUsecase.Login]: shop not found")
			return "", err
		}

		err = errors.Wrap(err, "[ShopUsecase.Login]: failed to get shop by name with password")
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(credentials.Password), []byte(password)); err != nil {
		err = errors.New("[ShopUsecase.Login]: invalid password")
		return "", err
	}

	t, err := utils.GenerateJWT(map[string]interface{}{
		"id":          credentials.ID,
		"name":        credentials.Name,
		"description": credentials.Description,
	}, []byte(viper.GetString("jwt.shopsecret")))

	if err != nil {
		err = errors.Wrap(err, "[ShopUsecase.Login]: failed to generate shop jwt")
		return "", err
	}

	return t, nil
}

func (u *shopUsecase) GetProductsByShopID(id uint32) ([]entity.Product, error) {
	log.Trace("Entering function GetProductsByShopID()")
	defer log.Trace("Exiting function GetProductsByShopID()")

	log.WithFields(log.Fields{
		"id": id,
	}).Debug("Getting products by shop id")

	// First check if shop exists
	exists, err := u.repo.ShopExists(id)
	if err != nil {
		err = errors.Wrap(err, "[ShopUsecase.GetProductsByShopID]: failed to check shop existence")
		return nil, err
	}
	if !exists {
		err = errors.New("[ShopUsecase.GetProductsByShopID]: shop not found")
		return nil, err
	}

	// If shop exists, get its products
	products, err := u.repo.GetProductsByShopID(id)
	if err != nil {
		err = errors.Wrap(err, "[ShopUsecase.GetProductsByShopID]: failed to get products by shop id")
		return nil, err
	}

	// Empty products is a valid case - return empty slice
	productsResponse := []entity.Product{}
	for _, product := range products {
		productsResponse = append(productsResponse, entity.Product{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}

	return productsResponse, nil
}

func (u *shopUsecase) UpdateProduct(req *entity.ProductManagementRequest, product *entity.Product) error {
	log.Trace("Entering function UpdateProduct()")
	defer log.Trace("Exiting function UpdateProduct()")

	log.WithFields(log.Fields{
		"req":     req,
		"product": product,
	}).Debug("Updating product")

	exists, err := u.repo.ShopExists(req.ShopID)
	if err != nil {
		err = errors.Wrap(err, "[ShopUsecase.UpdateProduct]: failed to check shop existence")
		return err
	}
	if !exists {
		err = errors.New("[ShopUsecase.UpdateProduct]: shop not found")
		return err
	}

	product.ID = req.ProductID
	if err := u.repo.UpdateProduct(req, product); err != nil {
		err = errors.Wrap(err, "[ShopUsecase.UpdateProduct]: failed to update product")
		return err
	}
	return nil
}

func (u *shopUsecase) DeleteProduct(req *entity.ProductManagementRequest) error {
	log.Trace("Entering function DeleteProduct()")
	defer log.Trace("Exiting function DeleteProduct()")

	log.WithFields(log.Fields{
		"req": req,
	}).Debug("Deleting product")

	exists, err := u.repo.ShopExists(req.ShopID)
	if err != nil {
		err = errors.Wrap(err, "[ShopUsecase.DeleteProduct]: failed to check shop existence")
		return err
	}
	if !exists {
		err = errors.New("[ShopUsecase.DeleteProduct]: shop not found")
		return err
	}

	if err := u.repo.DeleteProduct(req); err != nil {
		err = errors.Wrap(err, "[ShopUsecase.DeleteProduct]: failed to delete product")
		return err
	}
	return nil
}

package usecase

import (
	"fmt"
	"order-management/domain"
	"order-management/entity"
	"order-management/utils"

	"github.com/pkg/errors"
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
	if err := u.repo.CreateProduct(product, shopID); err != nil {
		return errors.Wrap(err, "[ShopUsecase.CreateProduct]: failed to create product")
	}
	return nil
}

func (u *shopUsecase) CreateShop(shop entity.Shop) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(shop.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, "[ShopUsecase.CreateShop]: failed to hash password")
	}

	shop.Password = string(hashedPassword)

	if err := u.repo.CreateShop(shop); err != nil {
		return errors.Wrap(err, "[ShopUsecase.CreateShop]: failed to create shop")
	}
	return nil
}

// Too big O(n^2)
func (u *shopUsecase) GetAllShopsWithProducts() ([]entity.ShopWithProducts, error) {
	shops, err := u.repo.GetAllShops()
	if err != nil {
		return nil, errors.Wrap(err, "[ShopUsecase.GetAllShopsWithProducts]: failed to get all shops")
	}
	fmt.Println("all shops was called :", len(shops))
	if len(shops) == 0 {
		return nil, errors.Wrap(gorm.ErrRecordNotFound, "[ShopUsecase.GetAllShopsWithProducts]: no shops found")
	}
	shopsResponse := []entity.ShopWithProducts{}
	for _, shop := range shops {
		products, err := u.repo.GetProductsByShopID(shop.ID)
		if err != nil {
			return nil, errors.Wrap(err, "[ShopUsecase.GetAllShopsWithProducts]: failed to get products by shop id")
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
	shops, err := u.repo.GetAllShops()
	if err != nil {
		return nil, errors.Wrap(err, "[ShopUsecase.GetAllShops]: failed to get all shops")
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
	shop, err := u.repo.GetShopByName(name)
	if err != nil {
		if err.Error() == "[ShopRepository.GetShopByName]: shop not found" {
			err = errors.New("[ShopUsecase.GetShopByName]: shop not found")
			return entity.ShopWithProducts{}, err
		}
		return entity.ShopWithProducts{}, errors.Wrap(err, "[ShopUsecase.GetShopByName]: failed to get shop by name")
	}

	products, err := u.repo.GetProductsByShopID(shop.ID)
	if err != nil {
		return entity.ShopWithProducts{}, errors.Wrap(err, "[ShopUsecase.GetShopByName]: failed to get products by shop id")
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
	credentials, err := u.repo.GetShopByNameWithPassword(name)
	if err != nil {
		if err.Error() == "[ShopRepository.GetShopByNameWithPassword]: shop not found" {
			return "", errors.New("[ShopUsecase.Login]: shop not found")
		}
		return "", errors.Wrap(err, "[ShopUsecase.Login]: failed to get shop by name with password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(credentials.Password), []byte(password)); err != nil {
		return "", errors.New("[ShopUsecase.Login]: invalid password")
	}

	t, err := utils.GenerateJWT(map[string]interface{}{
		"id":          credentials.ID,
		"name":        credentials.Name,
		"description": credentials.Description,
	}, []byte(viper.GetString("jwt.shopsecret")))

	if err != nil {
		return "", errors.Wrap(err, "[ShopUsecase.Login]: failed to generate shop jwt")
	}
	return t, nil
}

func (u *shopUsecase) GetProductsByShopID(id uint32) ([]entity.Product, error) {
	// First check if shop exists
	exists, err := u.repo.ShopExists(id)
	if err != nil {
		return nil, errors.Wrap(err, "[ShopUsecase.GetProductsByShopID]: failed to check shop existence")
	}
	if !exists {
		return nil, errors.New("[ShopUsecase.GetProductsByShopID]: shop not found")
	}

	// If shop exists, get its products
	products, err := u.repo.GetProductsByShopID(id)
	if err != nil {
		return nil, errors.Wrap(err, "[ShopUsecase.GetProductsByShopID]: failed to get products by shop id")
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
	exists, err := u.repo.ShopExists(req.ShopID)
	if err != nil {
		return errors.Wrap(err, "[ShopUsecase.UpdateProduct]: failed to check shop existence")
	}
	if !exists {
		return errors.New("[ShopUsecase.UpdateProduct]: shop not found")
	}

	product.ID = req.ProductID
	if err := u.repo.UpdateProduct(req, product); err != nil {
		return errors.Wrap(err, "[ShopUsecase.UpdateProduct]: failed to update product")
	}
	return nil
}
func (u *shopUsecase) DeleteProduct(req *entity.ProductManagementRequest) error {
	exists, err := u.repo.ShopExists(req.ShopID)
	if err != nil {
		return errors.Wrap(err, "[ShopUsecase.DeleteProduct]: failed to check shop existence")
	}
	if !exists {
		return errors.New("[ShopUsecase.DeleteProduct]: shop not found")
	}

	if err := u.repo.DeleteProduct(req); err != nil {
		return errors.Wrap(err, "[ShopUsecase.DeleteProduct]: failed to delete product")
	}
	return nil
}

package usecase

import (
	"fmt"
	"order-management/domain"
	"order-management/entity"
	"order-management/utils"

	"golang.org/x/crypto/bcrypt"
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(shop.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	shop.Password = string(hashedPassword)
	return u.repo.CreateShop(shop)
}

// Too big O(n^2)
func (u *shopUsecase) GetAllShopsWithProducts() ([]entity.ShopWithProducts, error) {
	shops, err := u.repo.GetAllShops()
	if err != nil {
		return nil, err
	}
	fmt.Println(shops)
	shopsResponse := []entity.ShopWithProducts{}
	for _, shop := range shops {
		products, err := u.repo.GetProductsByShopID(shop.ID)
		if err != nil {
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
	shops, err := u.repo.GetAllShops()
	if err != nil {
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
	shop, err := u.repo.GetShopByName(name)
	if err != nil {
		return entity.ShopWithProducts{}, err
	}
	products, err := u.repo.GetProductsByShopID(shop.ID)
	if err != nil {
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
	credentials, err := u.repo.GetShopByNameWithPassword(name)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(credentials.Password), []byte(password)); err != nil {
		return "", err
	}
	t, err := utils.GenerateShopJWT(&entity.ShopResponse{
		ID:          credentials.ID,
		Name:        credentials.Name,
		Description: credentials.Description,
	})
	if err != nil {
		return "", err
	}
	return t, nil
}

func (u *shopUsecase) GetProductsByShopID(id uint32) ([]entity.Product, error) {
	products, err := u.repo.GetProductsByShopID(id)
	if err != nil {
		return nil, err
	}
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

func (u *shopUsecase) BelongsToShop(productID uint32, claims *entity.ShopResponse) bool {
	product, err := u.repo.GetProductByID(productID)
	if err != nil {
		return false
	}
	return product.ShopID == claims.ID
}

func (u *shopUsecase) UpdateProduct(req *entity.ProductManagementRequest, newProduct *entity.Product) error {
	newProduct.ID = req.ProductID
	fmt.Println("new product ", newProduct)
	return u.repo.UpdateProduct(req, newProduct)
}
func (u *shopUsecase) DeleteProduct(req *entity.ProductManagementRequest) error {
	return u.repo.DeleteProduct(req)
}

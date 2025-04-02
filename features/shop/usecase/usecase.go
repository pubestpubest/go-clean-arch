package usecase

import (
	"fmt"
	"order-management/domain"
	"order-management/entity"
	"order-management/response"

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
func (u *shopUsecase) GetAllShopsWithProducts() ([]response.ShopWithProducts, error) {
	shops, err := u.repo.GetAllShops()
	if err != nil {
		return nil, err
	}
	fmt.Println(shops)
	shopsResponse := []response.ShopWithProducts{}
	for _, shop := range shops {
		products, err := u.repo.GetProductsByShopID(shop.ID)
		if err != nil {
			return nil, err
		}
		productsResponse := []response.Product{}
		for _, product := range products {
			productsResponse = append(productsResponse, response.Product{
				ID:          product.ID,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
			})
		}
		shopResponse := response.ShopWithProducts{
			ID:          shop.ID,
			Name:        shop.Name,
			Description: shop.Description,
			Products:    productsResponse,
		}
		shopsResponse = append(shopsResponse, shopResponse)
	}
	return shopsResponse, nil
}

func (u *shopUsecase) GetAllShops() ([]response.Shop, error) {
	shops, err := u.repo.GetAllShops()
	if err != nil {
		return nil, err
	}
	shopsResponse := []response.Shop{}
	for _, shop := range shops {
		shopsResponse = append(shopsResponse, response.Shop{
			ID:          shop.ID,
			Name:        shop.Name,
			Description: shop.Description,
		})
	}
	return shopsResponse, nil
}
func (u *shopUsecase) GetShopByName(name string) (entity.Shop, error) {
	return u.repo.GetShopByName(name)
}

func (u *shopUsecase) Login(name string, password string) (entity.Shop, error) {
	return u.repo.GetShopByName(name)
}
func (u *shopUsecase) GetProductsByShopID(id uint32) ([]response.Product, error) {
	products, err := u.repo.GetProductsByShopID(id)
	if err != nil {
		return nil, err
	}
	productsResponse := []response.Product{}
	for _, product := range products {
		productsResponse = append(productsResponse, response.Product{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}
	return productsResponse, nil
}

func (u *shopUsecase) BelongsToShop(productID uint32, claims *response.Shop) bool {
	product, err := u.repo.GetProductByID(productID)
	if err != nil {
		return false
	}
	return product.ShopID == claims.ID
}

func (u *shopUsecase) UpdateProduct(productID uint32, newProduct *entity.Product) error {
	newProduct.ID = productID
	fmt.Println("new product ", newProduct)
	return u.repo.UpdateProduct(productID, newProduct)
}

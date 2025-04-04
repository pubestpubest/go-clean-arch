package entity

type ProductManagementRequest struct {
	ShopWithOutPassword
	ProductID uint32 `json:"product_id"`
}

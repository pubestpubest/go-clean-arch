package entity

type ProductManagementRequest struct {
	ShopResponse
	ProductID uint32 `json:"product_id"`
}

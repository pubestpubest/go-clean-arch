package entity

type ProductManagementRequest struct {
	ShopID    uint32 `json:"shop_id"`
	ProductID uint32 `json:"product_id"`
}

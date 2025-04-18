package entity

type Product struct {
	ID            uint32 `gorm:"primary_key"`
	Name          string
	Description   string
	Price         uint32
	ShopID        uint32
	Shop          Shop           `gorm:"foreignKey:ShopID"`
	Orders        []Order        `gorm:"many2many:order_products;"`
	OrderProducts []OrderProduct `gorm:"foreignKey:ProductID"`
}

type ProductWithOutShop struct {
	ID          uint32 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       uint32 `json:"price"`
}

type ProductOrderAmount struct {
	ID          uint32 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       uint32 `json:"price"`
	Amount      uint32 `json:"amount"`
}

type ProductManagementRequest struct {
	ShopID    uint32 `json:"shop_id"`
	ProductID uint32 `json:"product_id"`
}

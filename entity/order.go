package entity

type Order struct {
	ID            uint32 `gorm:"primary_key"`
	Status        Status `gorm:"type:varchar(20)"`
	Total         float32
	Courier       string
	UserID        uint32
	User          User           `gorm:"foreignKey:UserID"`
	Products      []Product      `gorm:"many2many:order_products;"`
	OrderProducts []OrderProduct `gorm:"foreignKey:OrderID"`
}

type Status string

const (
	PENDING   Status = "PENDING"
	SHIPPING  Status = "SHIPPING"
	CANCELLED Status = "CANCELLED"
	COMPLETED Status = "COMPLETED"
)

// OrderProduct represents the join table between Order and Product with additional fields
type OrderProduct struct {
	OrderID   uint32  `gorm:"primaryKey"`
	ProductID uint32  `gorm:"primaryKey"`
	Amount    uint32  `gorm:"not null"` // Amount of products in the order
	Order     Order   `gorm:"foreignKey:OrderID"`
	Product   Product `gorm:"foreignKey:ProductID"`
}

type OrderRequest struct {
	OrderProducts []OrderProductRequest `json:"orderProducts"`
	Courier       string                `json:"courier"`
}

type OrderProductRequest struct {
	ProductId uint32 `json:"productId"`
	Amount    uint32 `json:"amount"`
}

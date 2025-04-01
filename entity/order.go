package entity

type Order struct {
	ID       uint32 `gorm:"primary_key"`
	Status   Status `gorm:"type:varchar(20)"`
	Total    float32
	Courier  string
	UserID   uint32
	User     User      `gorm:"foreignKey:UserID"`
	Products []Product `gorm:"many2many:order_products;"`
}

type Status string

const (
	PENDING   Status = "PENDING"
	SHIPPING  Status = "SHIPPING"
	CANCELLED Status = "CANCELLED"
	COMPLETED Status = "COMPLETED"
)

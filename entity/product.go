package entity

type Product struct {
	ID          uint32 `gorm:"primary_key"`
	Name        string
	Description string
	Price       uint32
	ShopID      uint32
	Shop        Shop    `gorm:"foreignKey:ShopID"`
	Orders      []Order `gorm:"many2many:order_products;"`
}

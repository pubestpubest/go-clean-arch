package entity

type Shop struct {
	ID          uint32 `gorm:"primary_key"`
	Name        string
	Description string
	Products    []Product `gorm:"foreignKey:ShopID"`
}

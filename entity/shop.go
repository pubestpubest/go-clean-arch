package entity

type Shop struct {
	ID          uint32 `gorm:"primary_key"`
	Name        string `gorm:"not null;unique"`
	Description string
	Password    string    `gorm:"not null"`
	Products    []Product `gorm:"foreignKey:ShopID"`
}

package entity

type User struct {
	ID      uint32 `gorm:"primary_key"`
	Email   string
	Address string
	Orders  []Order `gorm:"foreignKey:UserID"`
}

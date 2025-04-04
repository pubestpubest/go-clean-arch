package entity

type User struct {
	ID       uint32 `gorm:"primary_key"`
	Email    string `gorm:"unique;not null"`
	Address  string
	Password string  `gorm:"not null"`
	Orders   []Order `gorm:"foreignKey:UserID"`
}

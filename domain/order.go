package domain

import "order-management/entity"

type OrderUsecase interface {
	GetAllOrders() ([]entity.Order, error)
	GetOrder(orderID uint32) (entity.Order, error)
	GetOrdersByUserID(userID uint32) ([]entity.Order, error)
	GetOrdersByShopID(shopID uint32) ([]entity.Order, error)
	CreateOrder(orderRequest entity.OrderRequest, userID uint32) error
}

type OrderRepository interface {
	CreateOrder(order entity.Order) error
	UpdateOrder(order entity.Order) error
	DeleteOrder(orderID uint32) error
	GetOrder(orderID uint32) (entity.Order, error)
	GetOrdersByUserID(userID uint32) ([]entity.Order, error)
	GetOrdersByShopID(shopID uint32) ([]entity.Order, error)
	GetAllOrders() ([]entity.Order, error)
}

package domain

import "order-management/entity"

type OrderUsecase interface {
	GetAllOrders() ([]entity.Order, error)
	GetOrder(orderID uint32) (entity.OrderResponse, error)
	GetOrdersByUserID(userID uint32) ([]entity.OrderResponse, error)
	GetOrdersByShopID(shopID uint32) ([]entity.OrderResponse, error)
	CreateOrder(orderRequest entity.OrderRequest, userID uint32) error
}

type OrderRepository interface {
	CreateOrder(order entity.Order) error
	UpdateOrder(order entity.Order) error
	DeleteOrder(orderID uint32) error
	GetOrder(orderID uint32) (entity.Order, error)
	GetOrdersByUserID(userID uint32) ([]uint32, error)
	GetOrdersByShopID(shopID uint32) ([]uint32, error)
	GetAllOrders() ([]entity.Order, error)
}

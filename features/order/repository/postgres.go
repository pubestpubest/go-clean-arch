package repository

import (
	"order-management/domain"
	"order-management/entity"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) domain.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrder(order entity.Order) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create the order first
		if err := tx.Create(&order).Error; err != nil {
			return errors.Wrap(err, "[OrderRepository.CreateOrder]: failed to create order")
		}

		// No need to create order products again as they are already created with the order
		return nil
	})
}

func (r *orderRepository) GetOrder(orderID uint32) (entity.Order, error) {
	var order entity.Order
	if err := r.db.Preload("Products").Where("id = ?", orderID).First(&order).Error; err != nil {
		err = errors.Wrap(err, "[OrderRepository.GetOrder]: failed to get order")
		return entity.Order{}, err
	}
	return order, nil
}

func (r *orderRepository) GetOrdersByUserID(userID uint32) ([]uint32, error) {
	var orderIDs []uint32
	if err := r.db.Table("orders").Select("orders.id").Where("orders.user_id = ?", userID).Pluck("orders.id", &orderIDs).Error; err != nil {
		err = errors.Wrap(err, "[OrderRepository.GetOrdersByUserID]: failed to get orders by user id")
		return nil, err
	}
	return orderIDs, nil
}

func (r *orderRepository) GetOrdersByShopID(shopID uint32) ([]uint32, error) {
	var orderIDs []uint32
	if err := r.db.Table("orders").Select("orders.id").
		Joins("JOIN order_products op ON orders.id = op.order_id").
		Where("op.shop_id = ?", shopID).
		Pluck("orders.id", &orderIDs).Error; err != nil {
		err = errors.Wrap(err, "[OrderRepository.GetOrdersByShopID]: failed to get orders by shop id")
		return nil, err
	}
	return orderIDs, nil
}

func (r *orderRepository) UpdateOrder(order entity.Order) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Update order details
		if err := tx.Model(&entity.Order{}).Where("id = ?", order.ID).Updates(order).Error; err != nil {
			return errors.Wrap(err, "[OrderRepository.UpdateOrder]: failed to update order")
		}

		// Update order products
		if len(order.OrderProducts) > 0 {
			// Delete existing order products
			if err := tx.Where("order_id = ?", order.ID).Delete(&entity.OrderProduct{}).Error; err != nil {
				return errors.Wrap(err, "[OrderRepository.UpdateOrder]: failed to delete existing order products")
			}

			// Create new order products
			for _, orderProduct := range order.OrderProducts {
				orderProduct.OrderID = order.ID
				if err := tx.Create(&orderProduct).Error; err != nil {
					return errors.Wrap(err, "[OrderRepository.UpdateOrder]: failed to create new order product")
				}
			}
		}

		return nil
	})
}

func (r *orderRepository) DeleteOrder(orderID uint32) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete order products first (due to foreign key constraint)
		if err := tx.Where("order_id = ?", orderID).Delete(&entity.OrderProduct{}).Error; err != nil {
			return errors.Wrap(err, "[OrderRepository.DeleteOrder]: failed to delete order products")
		}

		// Delete the order
		if err := tx.Where("id = ?", orderID).Delete(&entity.Order{}).Error; err != nil {
			return errors.Wrap(err, "[OrderRepository.DeleteOrder]: failed to delete order")
		}

		return nil
	})
}

func (r *orderRepository) GetAllOrders() ([]entity.Order, error) {
	var orders []entity.Order
	if err := r.db.Preload("OrderProducts.Product").Find(&orders).Error; err != nil {
		return nil, errors.Wrap(err, "[OrderRepository.GetAllOrders]: failed to get all orders")
	}
	return orders, nil
}

func (r *orderRepository) GetProductOrderAmount(orderID uint32, productID uint32) (uint32, error) {
	var amount uint32
	if err := r.db.Table("order_products").Where("order_id = ? AND product_id = ?", orderID, productID).Pluck("amount", &amount).Error; err != nil {
		return 0, errors.Wrap(err, "[OrderRepository.GetProductOrderAmount]: failed to get product order amount")
	}
	return amount, nil
}

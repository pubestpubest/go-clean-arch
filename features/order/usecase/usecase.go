package usecase

import (
	"order-management/domain"
	"order-management/entity"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type OrderUsecase struct {
	orderRepo   domain.OrderRepository
	productRepo domain.ProductRepository
}

func NewOrderUsecase(orderRepo domain.OrderRepository, productRepo domain.ProductRepository) domain.OrderUsecase {
	return &OrderUsecase{orderRepo: orderRepo, productRepo: productRepo}
}

func (u *OrderUsecase) CreateOrder(orderRequest entity.OrderRequest, userID uint32) error {
	log.Trace("Entering function CreateOrder()")
	defer log.Trace("Exiting function CreateOrder()")

	log.WithFields(log.Fields{
		"orderRequest": orderRequest,
	}).Debug("Creating order")

	totalPrice := 0.0
	for _, reqProduct := range orderRequest.OrderProducts {
		price, err := u.productRepo.GetProductPrice(reqProduct.ProductId)
		if err != nil {
			err = errors.Wrap(err, "[OrderUsecase.CreateOrder]: failed to get product price")
			return err
		}
		totalPrice += price * float64(reqProduct.Amount)
	}
	// 1. Create the main order
	order := entity.Order{
		Status:        entity.PENDING, // Initial status should be PENDING
		Courier:       orderRequest.Courier,
		UserID:        userID,
		OrderProducts: make([]entity.OrderProduct, len(orderRequest.OrderProducts)),
		Total:         float32(totalPrice),
	}

	// 2. Transform OrderProductRequest into OrderProduct entries
	for i, reqProduct := range orderRequest.OrderProducts {
		orderProduct := entity.OrderProduct{
			ProductID: reqProduct.ProductId,
			Amount:    reqProduct.Amount,
			// OrderID will be automatically set by the repository after order creation
		}
		order.OrderProducts[i] = orderProduct
	}

	// 3. Call the repository to create the order
	if err := u.orderRepo.CreateOrder(order); err != nil {
		err = errors.Wrap(err, "[OrderUsecase.CreateOrder]: failed to create order")
		return err
	}

	return nil
}

func (u *OrderUsecase) GetAllOrders() ([]entity.Order, error) {
	log.Trace("Entering function GetAllOrders()")
	defer log.Trace("Exiting function GetAllOrders()")

	log.Debug("Getting all orders")

	orders, err := u.orderRepo.GetAllOrders()
	if err != nil {
		err = errors.Wrap(err, "[OrderUsecase.GetAllOrders]: failed to get all orders")
		return nil, err
	}

	return orders, nil
}

func (u *OrderUsecase) GetOrder(orderID uint32) (entity.Order, error) {
	log.Trace("Entering function GetOrder()")
	defer log.Trace("Exiting function GetOrder()")

	log.WithFields(log.Fields{
		"orderID": orderID,
	}).Debug("Getting order by ID")

	order, err := u.orderRepo.GetOrder(orderID)
	if err != nil {
		err = errors.Wrap(err, "[OrderUsecase.GetOrder]: failed to get order by ID")
		return entity.Order{}, err
	}
	return order, nil
}

func (u *OrderUsecase) GetOrdersByUserID(userID uint32) ([]entity.Order, error) {
	log.Trace("Entering function GetOrdersByUserID()")
	defer log.Trace("Exiting function GetOrdersByUserID()")

	log.WithFields(log.Fields{
		"userID": userID,
	}).Debug("Getting orders by user ID")

	orders, err := u.orderRepo.GetOrdersByUserID(userID)
	if err != nil {
		err = errors.Wrap(err, "[OrderUsecase.GetOrdersByUserID]: failed to get orders by user ID")
		return nil, err
	}
	return orders, nil
}

func (u *OrderUsecase) GetOrdersByShopID(shopID uint32) ([]entity.Order, error) {
	log.Trace("Entering function GetOrdersByShopID()")
	defer log.Trace("Exiting function GetOrdersByShopID()")

	log.WithFields(log.Fields{
		"shopID": shopID,
	}).Debug("Getting orders by shop ID")

	orders, err := u.orderRepo.GetOrdersByShopID(shopID)
	if err != nil {
		err = errors.Wrap(err, "[OrderUsecase.GetOrdersByShopID]: failed to get orders by shop ID")
		return nil, err
	}
	return orders, nil
}

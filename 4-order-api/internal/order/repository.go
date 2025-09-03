package order

import (
	"demo/order-api/internal/product"
	"demo/order-api/internal/user"
	"errors"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// CreateOrder создает новый заказ
func (r *OrderRepository) CreateOrder(order *Order) error {
	return r.db.Create(order).Error
}

// GetOrderByID получает заказ по ID с полной информацией
func (r *OrderRepository) GetOrderByID(id uint) (*Order, error) {
	var order Order
	err := r.db.Preload("User").
		Preload("OrderItems.Product").
		Where("id = ?", id).
		First(&order).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}
	
	return &order, nil
}

// GetOrdersByUserID получает все заказы пользователя
func (r *OrderRepository) GetOrdersByUserID(userID uint) ([]Order, error) {
	var orders []Order
	err := r.db.Preload("OrderItems.Product").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error
	
	if err != nil {
		return nil, err
	}
	
	return orders, nil
}

// GetUserByPhone получает пользователя по номеру телефона
func (r *OrderRepository) GetUserByPhone(phone string) (*user.User, error) {
	var user user.User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	
	return &user, nil
}

// GetProductsByIDs получает продукты по их ID
func (r *OrderRepository) GetProductsByIDs(ids []uint) ([]product.Product, error) {
	var products []product.Product
	err := r.db.Where("id IN ?", ids).Find(&products).Error
	
	if err != nil {
		return nil, err
	}
	
	return products, nil
}

// UpdateOrderStatus обновляет статус заказа
func (r *OrderRepository) UpdateOrderStatus(orderID uint, status OrderStatus) error {
	return r.db.Model(&Order{}).Where("id = ?", orderID).Update("status", status).Error
}

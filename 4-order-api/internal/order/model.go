package order

import (
	"demo/order-api/internal/product"
	"demo/order-api/internal/user"
	"time"

	"gorm.io/gorm"
)

// Order представляет заказ пользователя
type Order struct {
	gorm.Model
	UserID      uint        `gorm:"not null"`
	User        user.User   `gorm:"foreignKey:UserID"`
	Status      OrderStatus `gorm:"type:varchar(20);default:'pending'"`
	TotalAmount float64     `gorm:"type:decimal(10,2);not null"`
	OrderItems  []OrderItem `gorm:"foreignKey:OrderID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// OrderItem представляет элемент заказа (связь между заказом и продуктом)
type OrderItem struct {
	gorm.Model
	OrderID   uint           `gorm:"not null"`
	Order     Order          `gorm:"foreignKey:OrderID"`
	ProductID uint           `gorm:"not null"`
	Product   product.Product `gorm:"foreignKey:ProductID"`
	Quantity  int            `gorm:"not null;default:1"`
	Price     float64        `gorm:"type:decimal(10,2);not null"`
}

// OrderStatus представляет статус заказа
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// CreateOrderRequest представляет запрос на создание заказа
type CreateOrderRequest struct {
	ProductIDs []uint `json:"product_ids" validate:"required,min=1"`
	Quantities []int  `json:"quantities" validate:"required,min=1"`
}

// OrderResponse представляет ответ с информацией о заказе
type OrderResponse struct {
	ID          uint                `json:"id"`
	Status      OrderStatus         `json:"status"`
	TotalAmount float64             `json:"total_amount"`
	OrderItems  []OrderItemResponse `json:"order_items"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

// OrderItemResponse представляет ответ с информацией об элементе заказа
type OrderItemResponse struct {
	ID       uint                `json:"id"`
	Product  product.Product     `json:"product"`
	Quantity int                 `json:"quantity"`
	Price    float64             `json:"price"`
}

package user

import (
	"gorm.io/gorm"
)

// User представляет пользователя системы
type User struct {
	gorm.Model
	Phone  string `gorm:"uniqueIndex;not null"`
	Orders []Order `gorm:"foreignKey:UserID"`
}

// Order представляет заказ пользователя (для связи)
type Order struct {
	gorm.Model
	UserID      uint        `gorm:"not null"`
	User        User        `gorm:"foreignKey:UserID"`
	Status      string      `gorm:"type:varchar(20);default:'pending'"`
	TotalAmount float64     `gorm:"type:decimal(10,2);not null"`
	OrderItems  []OrderItem `gorm:"foreignKey:OrderID"`
}

// OrderItem представляет элемент заказа (для связи)
type OrderItem struct {
	gorm.Model
	OrderID   uint    `gorm:"not null"`
	Order     Order   `gorm:"foreignKey:OrderID"`
	ProductID uint    `gorm:"not null"`
	Quantity  int     `gorm:"not null;default:1"`
	Price     float64 `gorm:"type:decimal(10,2);not null"`
}

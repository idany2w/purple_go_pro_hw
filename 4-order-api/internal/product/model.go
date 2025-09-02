package product

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Product представляет продукт в системе
type Product struct {
	gorm.Model
	Name        string         `gorm:"not null"`
	Description string         `gorm:"type:text"`
	Images      pq.StringArray `gorm:"type:text[]"`
	OrderItems  []OrderItem    `gorm:"foreignKey:ProductID"`
}

// OrderItem представляет элемент заказа (для связи)
type OrderItem struct {
	gorm.Model
	OrderID   uint    `gorm:"not null"`
	ProductID uint    `gorm:"not null"`
	Product   Product `gorm:"foreignKey:ProductID"`
	Quantity  int     `gorm:"not null;default:1"`
	Price     float64 `gorm:"type:decimal(10,2);not null"`
}

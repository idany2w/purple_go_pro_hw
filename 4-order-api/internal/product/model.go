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
}

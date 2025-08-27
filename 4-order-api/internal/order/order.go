package order

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	Name        string         `gorm:"not null"`
	Description string         `gorm:"type:text"`
	Images      pq.StringArray `gorm:"type:text[]"`
}

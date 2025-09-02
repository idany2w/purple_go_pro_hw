package user

import (
	"gorm.io/gorm"
)

// User представляет пользователя системы
type User struct {
	gorm.Model
	Phone string `gorm:"uniqueIndex;not null"`
}

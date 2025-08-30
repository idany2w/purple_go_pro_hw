package sms

import "gorm.io/gorm"

type Sms struct {
	gorm.Model
	Phone string `gorm:"unique"`
	Token string `gorm:"unique"`
	Code  string `gorm:"unique"`
}

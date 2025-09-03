package user

import (
	"errors"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetUserByPhone получает пользователя по номеру телефона
func (r *UserRepository) GetUserByPhone(phone string) (*User, error) {
	var user User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	
	return &user, nil
}

// CreateUser создает нового пользователя
func (r *UserRepository) CreateUser(user *User) error {
	return r.db.Create(user).Error
}

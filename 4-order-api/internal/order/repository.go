package order

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (r *OrderRepository) GetAll(limit int, offset int) ([]Order, error) {
	var orders []Order
	result := r.DB.Limit(limit).Offset(offset).Find(&orders)

	if result.Error != nil {
		return nil, result.Error
	}

	return orders, nil
}

func (r *OrderRepository) GetById(id uint) (*Order, error) {
	order := &Order{}
	result := r.DB.First(order, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return order, nil
}

func (r *OrderRepository) Create(order *Order) (*Order, error) {
	result := r.DB.Create(order)

	if result.Error != nil {
		return nil, result.Error
	}

	return order, nil
}

func (r *OrderRepository) Update(order *Order) (*Order, error) {
	result := r.DB.Clauses(clause.Returning{}).Updates(order)

	if result.Error != nil {
		return nil, result.Error
	}

	return order, nil
}

func (r *OrderRepository) Delete(id uint) error {
	result := r.DB.Delete(&Order{}, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

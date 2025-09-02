package product

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductRepository struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (r *ProductRepository) GetAll(limit int, offset int) ([]Product, error) {
	var products []Product
	result := r.DB.Limit(limit).Offset(offset).Find(&products)

	if result.Error != nil {
		return nil, result.Error
	}

	return products, nil
}

func (r *ProductRepository) GetById(id uint) (*Product, error) {
	product := &Product{}
	result := r.DB.First(product, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return product, nil
}

func (r *ProductRepository) Create(product *Product) (*Product, error) {
	result := r.DB.Create(product)

	if result.Error != nil {
		return nil, result.Error
	}

	return product, nil
}

func (r *ProductRepository) Update(product *Product) (*Product, error) {
	result := r.DB.Clauses(clause.Returning{}).Updates(product)

	if result.Error != nil {
		return nil, result.Error
	}

	return product, nil
}

func (r *ProductRepository) Delete(id uint) error {
	result := r.DB.Delete(&Product{}, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// GetProductsByIDs получает продукты по их ID
func (r *ProductRepository) GetProductsByIDs(ids []uint) ([]Product, error) {
	var products []Product
	err := r.DB.Where("id IN ?", ids).Find(&products).Error
	
	if err != nil {
		return nil, err
	}
	
	return products, nil
}

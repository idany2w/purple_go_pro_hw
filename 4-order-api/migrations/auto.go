package migrations

import (
	"demo/order-api/internal/order"
	"demo/order-api/internal/product"
	"demo/order-api/internal/user"
	"log"

	"gorm.io/gorm"
)

// AutoMigrate выполняет автоматические миграции для всех моделей
func AutoMigrate(db *gorm.DB) error {
	log.Println("Starting database migrations...")

	// Миграция таблицы пользователей
	err := db.AutoMigrate(&user.User{})
	if err != nil {
		log.Printf("Failed to migrate users table: %v", err)
		return err
	}
	log.Println("Users table migrated successfully")

	// Миграция таблицы продуктов
	err = db.AutoMigrate(&product.Product{})
	if err != nil {
		log.Printf("Failed to migrate products table: %v", err)
		return err
	}
	log.Println("Products table migrated successfully")

	// Миграция таблицы заказов
	err = db.AutoMigrate(&order.Order{})
	if err != nil {
		log.Printf("Failed to migrate orders table: %v", err)
		return err
	}
	log.Println("Orders table migrated successfully")

	// Миграция таблицы элементов заказа
	err = db.AutoMigrate(&order.OrderItem{})
	if err != nil {
		log.Printf("Failed to migrate order_items table: %v", err)
		return err
	}
	log.Println("Order items table migrated successfully")

	log.Println("All database migrations completed successfully")
	return nil
}

// CreateSampleData создает тестовые данные для демонстрации
func CreateSampleData(db *gorm.DB) error {
	log.Println("Creating sample data...")

	// Проверяем, есть ли уже данные
	var userCount int64
	db.Model(&user.User{}).Count(&userCount)
	if userCount > 0 {
		log.Println("Sample data already exists, skipping...")
		return nil
	}

	// Создаем тестового пользователя
	testUser := &user.User{
		Phone: "+79001234567",
	}
	err := db.Create(testUser).Error
	if err != nil {
		log.Printf("Failed to create test user: %v", err)
		return err
	}

	// Создаем тестовые продукты
	testProducts := []product.Product{
		{
			Name:        "iPhone 15",
			Description: "Смартфон Apple iPhone 15",
			Images:      []string{"iphone15_1.jpg", "iphone15_2.jpg"},
		},
		{
			Name:        "MacBook Pro",
			Description: "Ноутбук Apple MacBook Pro",
			Images:      []string{"macbook_pro_1.jpg", "macbook_pro_2.jpg"},
		},
		{
			Name:        "AirPods Pro",
			Description: "Беспроводные наушники Apple AirPods Pro",
			Images:      []string{"airpods_pro_1.jpg"},
		},
	}

	for _, product := range testProducts {
		err := db.Create(&product).Error
		if err != nil {
			log.Printf("Failed to create test product %s: %v", product.Name, err)
			return err
		}
	}

	log.Println("Sample data created successfully")
	return nil
}

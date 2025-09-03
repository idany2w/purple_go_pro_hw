package order

import (
	"demo/order-api/internal/product"
	"demo/order-api/internal/user"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// MockOrderRepository мок для OrderRepository
type MockOrderRepository struct {
	users    map[string]*user.User
	products map[uint]*product.Product
	orders   []*Order
}

func NewMockOrderRepository() *MockOrderRepository {
	return &MockOrderRepository{
		users:    make(map[string]*user.User),
		products: make(map[uint]*product.Product),
		orders:   make([]*Order, 0),
	}
}

func (m *MockOrderRepository) GetUserByPhone(phone string) (*user.User, error) {
	if user, exists := m.users[phone]; exists {
		return user, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *MockOrderRepository) GetProductsByIDs(ids []uint) ([]product.Product, error) {
	var products []product.Product
	for _, id := range ids {
		if product, exists := m.products[id]; exists {
			products = append(products, *product)
		}
	}
	return products, nil
}

func (m *MockOrderRepository) CreateOrder(order *Order) error {
	order.ID = uint(len(m.orders) + 1)
	m.orders = append(m.orders, order)
	return nil
}

func (m *MockOrderRepository) GetOrderByID(orderID uint) (*Order, error) {
	for _, order := range m.orders {
		if order.ID == orderID {
			return order, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *MockOrderRepository) GetOrdersByUserID(userID uint) ([]Order, error) {
	var userOrders []Order
	for _, order := range m.orders {
		if order.UserID == userID {
			userOrders = append(userOrders, *order)
		}
	}
	return userOrders, nil
}

// TestCreateOrderSuccess тестирует успешное создание заказа
func TestCreateOrderSuccess(t *testing.T) {
	// Создаем мок репозиторий
	mockRepo := NewMockOrderRepository()

	// Добавляем тестового пользователя
	testUser := &user.User{
		Model: gorm.Model{ID: 1},
		Phone: "+79001234567",
	}
	mockRepo.users[testUser.Phone] = testUser

	// Добавляем тестовый продукт
	testProduct := &product.Product{
		Model:       gorm.Model{ID: 1},
		Name:        "Test iPhone",
		Description: "Test smartphone",
		Images:      []string{"test1.jpg", "test2.jpg"},
	}
	mockRepo.products[testProduct.ID] = testProduct

	// Создаем сервис
	orderService := NewOrderService(mockRepo)

	// Создаем запрос на создание заказа
	createOrderRequest := &CreateOrderRequest{
		ProductIDs: []uint{testProduct.ID},
		Quantities: []int{2},
	}

	// Выполняем создание заказа
	order, err := orderService.CreateOrder(testUser.Phone, createOrderRequest)

	// Проверяем результат
	require.NoError(t, err, "Expected no error")
	assert.NotNil(t, order, "Expected order to be created")
	assert.Equal(t, testUser.ID, order.UserID, "Expected correct user ID")
	assert.Equal(t, OrderStatusPending, order.Status, "Expected pending status")
	assert.Equal(t, 200.0, order.TotalAmount, "Expected correct total amount")
	assert.Len(t, order.OrderItems, 1, "Expected one order item")

	// Проверяем элемент заказа
	orderItem := order.OrderItems[0]
	assert.Equal(t, testProduct.ID, orderItem.ProductID, "Expected correct product ID")
	assert.Equal(t, 2, orderItem.Quantity, "Expected correct quantity")
	assert.Equal(t, 100.0, orderItem.Price, "Expected correct price")
}

// TestCreateOrderWithoutProducts тестирует создание заказа без продуктов
func TestCreateOrderWithoutProducts(t *testing.T) {
	mockRepo := NewMockOrderRepository()
	orderService := NewOrderService(mockRepo)

	// Создаем запрос без продуктов
	createOrderRequest := &CreateOrderRequest{
		ProductIDs: []uint{},
		Quantities: []int{},
	}

	// Выполняем создание заказа
	order, err := orderService.CreateOrder("+79001234567", createOrderRequest)

	// Проверяем, что получили ошибку
	assert.Error(t, err, "Expected error")
	assert.Nil(t, order, "Expected no order")
	assert.Contains(t, err.Error(), "at least one product is required", "Expected specific error message")
}

// TestCreateOrderWithMismatchedArrays тестирует создание заказа с несоответствующими массивами
func TestCreateOrderWithMismatchedArrays(t *testing.T) {
	mockRepo := NewMockOrderRepository()
	orderService := NewOrderService(mockRepo)

	// Создаем запрос с несоответствующими массивами
	createOrderRequest := &CreateOrderRequest{
		ProductIDs: []uint{1, 2},
		Quantities: []int{1}, // Меньше чем productIDs
	}

	// Выполняем создание заказа
	order, err := orderService.CreateOrder("+79001234567", createOrderRequest)

	// Проверяем, что получили ошибку
	assert.Error(t, err, "Expected error")
	assert.Nil(t, order, "Expected no order")
	assert.Contains(t, err.Error(), "must have the same length", "Expected specific error message")
}

// TestCreateOrderWithInvalidQuantity тестирует создание заказа с невалидным количеством
func TestCreateOrderWithInvalidQuantity(t *testing.T) {
	// Создаем мок репозиторий
	mockRepo := NewMockOrderRepository()

	// Добавляем тестового пользователя
	testUser := &user.User{
		Model: gorm.Model{ID: 1},
		Phone: "+79001234567",
	}
	mockRepo.users[testUser.Phone] = testUser

	// Добавляем тестовый продукт
	testProduct := &product.Product{
		Model:       gorm.Model{ID: 1},
		Name:        "Test iPhone",
		Description: "Test smartphone",
		Images:      []string{"test1.jpg", "test2.jpg"},
	}
	mockRepo.products[testProduct.ID] = testProduct

	// Создаем сервис
	orderService := NewOrderService(mockRepo)

	// Создаем запрос с невалидным количеством
	createOrderRequest := &CreateOrderRequest{
		ProductIDs: []uint{testProduct.ID},
		Quantities: []int{0}, // Невалидное количество
	}

	// Выполняем создание заказа
	order, err := orderService.CreateOrder(testUser.Phone, createOrderRequest)

	// Проверяем, что получили ошибку
	assert.Error(t, err, "Expected error")
	assert.Nil(t, order, "Expected no order")
	assert.Contains(t, err.Error(), "invalid quantity", "Expected specific error message")
}

// TestCreateOrderWithNonExistentUser тестирует создание заказа с несуществующим пользователем
func TestCreateOrderWithNonExistentUser(t *testing.T) {
	mockRepo := NewMockOrderRepository()
	orderService := NewOrderService(mockRepo)

	// Создаем запрос
	createOrderRequest := &CreateOrderRequest{
		ProductIDs: []uint{1},
		Quantities: []int{1},
	}

	// Выполняем создание заказа с несуществующим пользователем
	order, err := orderService.CreateOrder("+79001234567", createOrderRequest)

	// Проверяем, что получили ошибку
	assert.Error(t, err, "Expected error")
	assert.Nil(t, order, "Expected no order")
	assert.Contains(t, err.Error(), "failed to get user", "Expected specific error message")
}

// TestCreateOrderWithNonExistentProduct тестирует создание заказа с несуществующим продуктом
func TestCreateOrderWithNonExistentProduct(t *testing.T) {
	// Создаем мок репозиторий
	mockRepo := NewMockOrderRepository()

	// Добавляем тестового пользователя
	testUser := &user.User{
		Model: gorm.Model{ID: 1},
		Phone: "+79001234567",
	}
	mockRepo.users[testUser.Phone] = testUser

	// Создаем сервис
	orderService := NewOrderService(mockRepo)

	// Создаем запрос с несуществующим продуктом
	createOrderRequest := &CreateOrderRequest{
		ProductIDs: []uint{99999}, // Несуществующий ID
		Quantities: []int{1},
	}

	// Выполняем создание заказа
	order, err := orderService.CreateOrder(testUser.Phone, createOrderRequest)

	// Проверяем, что получили ошибку
	assert.Error(t, err, "Expected error")
	assert.Nil(t, order, "Expected no order")
	assert.Contains(t, err.Error(), "some products not found", "Expected specific error message")
}

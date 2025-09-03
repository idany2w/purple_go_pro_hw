package order_test

import (
	"demo/order-api/internal/order"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func bootstrap() (*order.OrderService, sqlmock.Sqlmock, error) {
	database, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create mock database: %v", err)
	}

	gormDb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: database,
	}), &gorm.Config{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create gorm database: %v", err)
	}

	// Создаем репозитории
	orderRepository := order.NewOrderRepository(gormDb)

	// Создаем сервисы
	orderService := order.NewOrderService(orderRepository)

	return orderService, mock, nil
}

func TestCreateOrderSuccess(t *testing.T) {
	orderService, mock, err := bootstrap()
	if err != nil {
		t.Fatalf("Failed to bootstrap: %v", err)
		return
	}

	// Ожидаем запрос на получение пользователя
	mock.ExpectQuery("SELECT (.+) FROM \"users\"").
		WillReturnRows(sqlmock.NewRows([]string{"id", "phone"}).
			AddRow(1, "+79001234567"))

	// Ожидаем запрос на получение продукта
	mock.ExpectQuery("SELECT (.+) FROM \"products\"").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "images"}).
			AddRow(1, "Test iPhone", "Test smartphone", "{test1.jpg,test2.jpg}"))

	// Ожидаем создание заказа и элементов заказа в одной транзакции
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO \"orders\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery("INSERT INTO \"order_items\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	// Создаем запрос на создание заказа
	createOrderRequest := &order.CreateOrderRequest{
		ProductIDs: []uint{1},
		Quantities: []int{2},
	}

	// Выполняем создание заказа
	order, err := orderService.CreateOrder("+79001234567", createOrderRequest)

	// Проверяем результат
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if order == nil {
		t.Fatalf("Expected order to be created")
	}

	if order.UserID != 1 {
		t.Fatalf("Expected user ID 1, got %d", order.UserID)
	}

	if order.Status != "pending" {
		t.Fatalf("Expected status pending, got %s", order.Status)
	}

	if order.TotalAmount != 200.0 {
		t.Fatalf("Expected total amount 200.0, got %f", order.TotalAmount)
	}

	if len(order.OrderItems) != 1 {
		t.Fatalf("Expected 1 order item, got %d", len(order.OrderItems))
	}

	// Проверяем, что все ожидаемые SQL запросы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Not all SQL expectations were met: %v", err)
	}
}

func TestCreateOrderWithoutProducts(t *testing.T) {
	orderService, _, err := bootstrap()
	if err != nil {
		t.Fatalf("Failed to bootstrap: %v", err)
		return
	}

	// Создаем запрос без продуктов
	createOrderRequest := &order.CreateOrderRequest{
		ProductIDs: []uint{},
		Quantities: []int{},
	}

	// Выполняем создание заказа
	order, err := orderService.CreateOrder("+79001234567", createOrderRequest)

	// Проверяем, что получили ошибку
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}

	if order != nil {
		t.Fatalf("Expected no order, got %v", order)
	}
}

func TestCreateOrderWithMismatchedArrays(t *testing.T) {
	orderService, _, err := bootstrap()
	if err != nil {
		t.Fatalf("Failed to bootstrap: %v", err)
		return
	}

	// Создаем запрос с несоответствующими массивами
	createOrderRequest := &order.CreateOrderRequest{
		ProductIDs: []uint{1, 2},
		Quantities: []int{1}, // Меньше чем productIDs
	}

	// Выполняем создание заказа
	order, err := orderService.CreateOrder("+79001234567", createOrderRequest)

	// Проверяем, что получили ошибку
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}

	if order != nil {
		t.Fatalf("Expected no order, got %v", order)
	}
}

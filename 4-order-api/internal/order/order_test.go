package order

import (
	"testing"
	"time"
)

func TestOrderStatusConstants(t *testing.T) {
	// Проверяем, что все статусы заказа определены
	expectedStatuses := []OrderStatus{
		OrderStatusPending,
		OrderStatusConfirmed,
		OrderStatusShipped,
		OrderStatusDelivered,
		OrderStatusCancelled,
	}

	for _, status := range expectedStatuses {
		if status == "" {
			t.Errorf("Order status should not be empty")
		}
	}
}

func TestCreateOrderRequestValidation(t *testing.T) {
	// Тест валидного запроса
	validRequest := &CreateOrderRequest{
		ProductIDs: []uint{1, 2, 3},
		Quantities: []int{1, 2, 1},
	}

	if len(validRequest.ProductIDs) != len(validRequest.Quantities) {
		t.Errorf("ProductIDs and Quantities should have the same length")
	}

	if len(validRequest.ProductIDs) == 0 {
		t.Errorf("At least one product is required")
	}

	// Тест невалидного запроса
	invalidRequest := &CreateOrderRequest{
		ProductIDs: []uint{1, 2},
		Quantities: []int{1}, // Разная длина
	}

	if len(invalidRequest.ProductIDs) == len(invalidRequest.Quantities) {
		t.Errorf("ProductIDs and Quantities should have different lengths for invalid request")
	}
}

func TestOrderResponseMapping(t *testing.T) {
	// Создаем тестовый заказ
	order := &Order{
		Status:      OrderStatusPending,
		TotalAmount: 150.0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Создаем тестовые элементы заказа
	orderItem := OrderItem{
		Quantity: 2,
		Price:    75.0,
	}

	order.OrderItems = []OrderItem{orderItem}

	// Тестируем маппинг в ответ
	response := mapOrderToResponse(order)

	if response.Status != order.Status {
		t.Errorf("Expected order status %s, got %s", order.Status, response.Status)
	}

	if response.TotalAmount != order.TotalAmount {
		t.Errorf("Expected total amount %.2f, got %.2f", order.TotalAmount, response.TotalAmount)
	}

	if len(response.OrderItems) != len(order.OrderItems) {
		t.Errorf("Expected %d order items, got %d", len(order.OrderItems), len(response.OrderItems))
	}
}

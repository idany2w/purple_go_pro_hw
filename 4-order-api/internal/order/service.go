package order

import (
	"demo/order-api/internal/product"
	"errors"
	"fmt"
)

type OrderService struct {
	orderRepository *OrderRepository
}

func NewOrderService(orderRepository *OrderRepository) *OrderService {
	return &OrderService{
		orderRepository: orderRepository,
	}
}

// CreateOrder создает новый заказ
func (s *OrderService) CreateOrder(userPhone string, request *CreateOrderRequest) (*Order, error) {
	// Валидация входных данных
	if len(request.ProductIDs) != len(request.Quantities) {
		return nil, errors.New("product_ids and quantities arrays must have the same length")
	}
	
	if len(request.ProductIDs) == 0 {
		return nil, errors.New("at least one product is required")
	}
	
	// Получаем пользователя
	user, err := s.orderRepository.GetUserByPhone(userPhone)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	// Получаем продукты
	products, err := s.orderRepository.GetProductsByIDs(request.ProductIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	
	if len(products) != len(request.ProductIDs) {
		return nil, errors.New("some products not found")
	}
	
	// Создаем заказ
	order := &Order{
		UserID:      user.ID,
		Status:      OrderStatusPending,
		TotalAmount: 0,
	}
	
	// Создаем элементы заказа
	var totalAmount float64
	for i, productID := range request.ProductIDs {
		quantity := request.Quantities[i]
		if quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for product %d", productID)
		}
		
		// Находим продукт для получения цены
		var product *product.Product
		for _, p := range products {
			if p.ID == productID {
				product = &p
				break
			}
		}
		
		if product == nil {
			return nil, fmt.Errorf("product %d not found", productID)
		}
		
		// Здесь должна быть логика получения цены продукта
		// Для демонстрации используем фиксированную цену
		price := 100.0 // В реальном приложении цена должна браться из продукта
		
		orderItem := OrderItem{
			ProductID: productID,
			Quantity:  quantity,
			Price:     price,
		}
		
		order.OrderItems = append(order.OrderItems, orderItem)
		totalAmount += price * float64(quantity)
	}
	
	order.TotalAmount = totalAmount
	
	// Сохраняем заказ в базе данных
	err = s.orderRepository.CreateOrder(order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	
	return order, nil
}

// GetOrderByID получает заказ по ID
func (s *OrderService) GetOrderByID(orderID uint) (*Order, error) {
	return s.orderRepository.GetOrderByID(orderID)
}

// GetOrdersByUser получает заказы пользователя по номеру телефона
func (s *OrderService) GetOrdersByUser(userPhone string) ([]Order, error) {
	user, err := s.orderRepository.GetUserByPhone(userPhone)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	return s.orderRepository.GetOrdersByUserID(user.ID)
}

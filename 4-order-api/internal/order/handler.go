package order

import (
	"demo/order-api/pkg/jwt"
	"demo/order-api/pkg/request"
	"demo/order-api/pkg/response"
	"net/http"
	"strconv"
	"strings"
)

type OrderHandlerDeps struct {
	OrderService *OrderService
	JWT          *jwt.JWT
}

type OrderHandler struct {
	orderService *OrderService
	jwt          *jwt.JWT
}

func NewOrderHandler(router *http.ServeMux, deps OrderHandlerDeps) {
	orderHandler := &OrderHandler{
		orderService: deps.OrderService,
		jwt:          deps.JWT,
	}

	router.HandleFunc("POST /order", orderHandler.createOrder())
	router.HandleFunc("GET /order/{id}", orderHandler.getOrderByID())
	router.HandleFunc("GET /my-orders", orderHandler.getMyOrders())
}

// createOrder создает новый заказ
func (h *OrderHandler) createOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем токен из заголовка
		token := extractTokenFromHeader(r)
		if token == "" {
			response.SendJsonError(&w, "Authorization token required", http.StatusUnauthorized)
			return
		}

		// Валидируем токен и получаем данные пользователя
		valid, claims := h.jwt.Parse(token)
		if !valid || claims == nil {
			response.SendJsonError(&w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Парсим тело запроса
		body, _ := request.HandleBody[CreateOrderRequest](w, r)
		if body == nil {
			response.SendJsonError(&w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Создаем заказ
		order, err := h.orderService.CreateOrder(claims.Phone, body)
		if err != nil {
			response.SendJsonError(&w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Формируем ответ
		orderResponse := mapOrderToResponse(order)
		response.SendJsonSuccess(&w, "Order created successfully", orderResponse, http.StatusCreated)
	}
}

// getOrderByID получает заказ по ID
func (h *OrderHandler) getOrderByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем ID заказа из URL
		path := strings.TrimPrefix(r.URL.Path, "/order/")
		orderID, err := strconv.ParseUint(path, 10, 32)
		if err != nil {
			response.SendJsonError(&w, "Invalid order ID", http.StatusBadRequest)
			return
		}

		// Получаем заказ
		order, err := h.orderService.GetOrderByID(uint(orderID))
		if err != nil {
			if err.Error() == "order not found" {
				response.SendJsonError(&w, "Order not found", http.StatusNotFound)
				return
			}
			response.SendJsonError(&w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Формируем ответ
		orderResponse := mapOrderToResponse(order)
		response.SendJsonSuccess(&w, "Order retrieved successfully", orderResponse, http.StatusOK)
	}
}

// getMyOrders получает заказы текущего пользователя
func (h *OrderHandler) getMyOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем токен из заголовка
		token := extractTokenFromHeader(r)
		if token == "" {
			response.SendJsonError(&w, "Authorization token required", http.StatusUnauthorized)
			return
		}

		// Валидируем токен и получаем данные пользователя
		valid, claims := h.jwt.Parse(token)
		if !valid || claims == nil {
			response.SendJsonError(&w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Получаем заказы пользователя
		orders, err := h.orderService.GetOrdersByUser(claims.Phone)
		if err != nil {
			response.SendJsonError(&w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Формируем ответ
		var orderResponses []OrderResponse
		for _, order := range orders {
			orderResponses = append(orderResponses, *mapOrderToResponse(&order))
		}

		response.SendJsonSuccess(&w, "Orders retrieved successfully", orderResponses, http.StatusOK)
	}
}

// extractTokenFromHeader извлекает токен из заголовка Authorization
func extractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	// Ожидаем формат: "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

// mapOrderToResponse преобразует модель заказа в ответ API
func mapOrderToResponse(order *Order) *OrderResponse {
	var orderItemResponses []OrderItemResponse

	for _, item := range order.OrderItems {
		orderItemResponse := OrderItemResponse{
			ID:       item.ID,
			Product:  item.Product,
			Quantity: item.Quantity,
			Price:    item.Price,
		}
		orderItemResponses = append(orderItemResponses, orderItemResponse)
	}

	return &OrderResponse{
		ID:          order.ID,
		Status:      order.Status,
		TotalAmount: order.TotalAmount,
		OrderItems:  orderItemResponses,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}
}

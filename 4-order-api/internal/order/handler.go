package order

import (
	"demo/order-api/configs"
	"demo/order-api/pkg/request"
	"demo/order-api/pkg/response"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type OrderHandlerDeps struct {
	*configs.Config
	OrderRepository *OrderRepository
}

type OrderHandler struct {
	*configs.Config
	OrderRepository *OrderRepository
}

func NewOrderHandler(router *http.ServeMux, deps OrderHandlerDeps) {
	orderHandler := &OrderHandler{
		Config:          deps.Config,
		OrderRepository: deps.OrderRepository,
	}

	router.HandleFunc("GET /products", orderHandler.index())
	router.HandleFunc("GET /products/{id}", orderHandler.get())

	router.HandleFunc("POST /products", orderHandler.create())
	router.HandleFunc("PATCH /products/{id}", orderHandler.update())
	router.HandleFunc("DELETE /products/{id}", orderHandler.delete())
}

func (h *OrderHandler) index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))

		if err != nil || limit < 1 || limit > 100 {
			limit = 10
		}

		page, err := strconv.Atoi(r.URL.Query().Get("page"))

		if err != nil || page < 1 || page > 100 {
			page = 1
		}

		offset := (page - 1) * limit

		orders, err := h.OrderRepository.GetAll(limit, offset)

		if err != nil {
			response.SendJsonError(&w, "Failed to fetch orders", http.StatusInternalServerError)
			return
		}

		response.SendJsonSuccess(&w, "Orders fetched successfully", orders, http.StatusOK)
	}
}

func (h *OrderHandler) get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getUintId(w, r)

		if err != nil {
			return
		}

		order, err := h.OrderRepository.GetById(id)

		if err != nil && err != gorm.ErrRecordNotFound {
			response.SendJsonError(&w, "Failed to fetch order", http.StatusInternalServerError)
			return
		}

		if order == nil {
			response.SendJsonError(&w, "Order not found", http.StatusNotFound)
			return
		}

		response.SendJsonSuccess(&w, "Order fetched successfully", order, http.StatusOK)
	}
}

func (h *OrderHandler) create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := request.HandleBody[CreateOrderPayload](w, r)

		if body == nil {
			return
		}

		order, err := h.OrderRepository.Create(&Order{
			Name:        body.Name,
			Description: body.Description,
			Images:      body.Images,
		})

		if err != nil {
			response.SendJsonError(&w, "Failed to create order", http.StatusInternalServerError)
			return
		}

		response.SendJsonSuccess(&w, "Order created successfully", order, http.StatusOK)
	}
}

func (h *OrderHandler) update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := request.HandleBody[UpdateOrderPayload](w, r)

		if body == nil {
			return
		}

		id, err := getUintId(w, r)

		if err != nil {
			return
		}

		order, err := h.OrderRepository.GetById(id)
		if err != nil && err != gorm.ErrRecordNotFound {
			response.SendJsonError(&w, "Failed to update order", http.StatusInternalServerError)
			return
		}

		if order == nil {
			response.SendJsonError(&w, "Order not found", http.StatusNotFound)
			return
		}

		order, err = h.OrderRepository.Update(&Order{
			Model:       gorm.Model{ID: id},
			Name:        body.Name,
			Description: body.Description,
			Images:      body.Images,
		})

		if err != nil {
			response.SendJsonError(&w, "Failed to update order", http.StatusInternalServerError)
			return
		}

		response.SendJsonSuccess(&w, "Order updated successfully", order, http.StatusOK)
	}
}

func (h *OrderHandler) delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getUintId(w, r)

		if err != nil {
			return
		}

		order, err := h.OrderRepository.GetById(id)

		if err != nil && err != gorm.ErrRecordNotFound {
			response.SendJsonError(&w, "Failed to delete order", http.StatusInternalServerError)
			return
		}

		if order == nil {
			response.SendJsonError(&w, "Order not found", http.StatusNotFound)
			return
		}

		err = h.OrderRepository.Delete(order.ID)

		if err != nil {
			response.SendJsonError(&w, "Failed to delete order", http.StatusInternalServerError)
			return
		}

		response.SendJsonSuccess(&w, "Order deleted successfully", order, http.StatusOK)
	}
}

func getUintId(w http.ResponseWriter, r *http.Request) (uint, error) {
	idString := r.PathValue("id")
	id, err := strconv.ParseUint(idString, 10, 32)

	if err != nil {
		response.SendJsonError(&w, "Invalid order id", http.StatusBadRequest)
		return 0, err
	}

	return uint(id), nil
}

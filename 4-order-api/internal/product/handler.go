package product

import (
	"demo/order-api/configs"
	"demo/order-api/pkg/request"
	"demo/order-api/pkg/response"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type ProductHandlerDeps struct {
	*configs.Config
	ProductRepository *ProductRepository
}

type ProductHandler struct {
	*configs.Config
	ProductRepository *ProductRepository
}

func NewProductHandler(router *http.ServeMux, deps ProductHandlerDeps) {
	productHandler := &ProductHandler{
		Config:            deps.Config,
		ProductRepository: deps.ProductRepository,
	}

	router.HandleFunc("GET /products", productHandler.index())
	router.HandleFunc("GET /products/{id}", productHandler.get())

	router.HandleFunc("POST /products", productHandler.create())
	router.HandleFunc("PATCH /products/{id}", productHandler.update())
	router.HandleFunc("DELETE /products/{id}", productHandler.delete())
}

func (h *ProductHandler) index() http.HandlerFunc {
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

		products, err := h.ProductRepository.GetAll(limit, offset)

		if err != nil {
			response.SendJsonError(&w, "Failed to fetch orders", http.StatusInternalServerError)
			return
		}

		response.SendJsonSuccess(&w, "Products fetched successfully", products, http.StatusOK)
	}
}

func (h *ProductHandler) get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getUintId(w, r)

		if err != nil {
			return
		}

		product, err := h.ProductRepository.GetById(id)

		if err != nil && err != gorm.ErrRecordNotFound {
			response.SendJsonError(&w, "Failed to fetch order", http.StatusInternalServerError)
			return
		}

		if product == nil {
			response.SendJsonError(&w, "Product not found", http.StatusNotFound)
			return
		}

		response.SendJsonSuccess(&w, "Product fetched successfully", product, http.StatusOK)
	}
}

func (h *ProductHandler) create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := request.HandleBody[CreateProductPayload](w, r)

		if body == nil {
			return
		}

		product, err := h.ProductRepository.Create(&Product{
			Name:        body.Name,
			Description: body.Description,
			Images:      body.Images,
		})

		if err != nil {
			response.SendJsonError(&w, "Failed to create product", http.StatusInternalServerError)
			return
		}

		response.SendJsonSuccess(&w, "Product created successfully", product, http.StatusOK)
	}
}

func (h *ProductHandler) update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := request.HandleBody[UpdateProductPayload](w, r)

		if body == nil {
			return
		}

		id, err := getUintId(w, r)

		if err != nil {
			return
		}

		product, err := h.ProductRepository.GetById(id)
		if err != nil && err != gorm.ErrRecordNotFound {
			response.SendJsonError(&w, "Failed to update order", http.StatusInternalServerError)
			return
		}

		if product == nil {
			response.SendJsonError(&w, "Product not found", http.StatusNotFound)
			return
		}

		product, err = h.ProductRepository.Update(&Product{
			Model:       gorm.Model{ID: id},
			Name:        body.Name,
			Description: body.Description,
			Images:      body.Images,
		})

		if err != nil {
			response.SendJsonError(&w, "Failed to update product", http.StatusInternalServerError)
			return
		}

		response.SendJsonSuccess(&w, "Product updated successfully", product, http.StatusOK)
	}
}

func (h *ProductHandler) delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getUintId(w, r)

		if err != nil {
			return
		}

		product, err := h.ProductRepository.GetById(id)

		if err != nil && err != gorm.ErrRecordNotFound {
			response.SendJsonError(&w, "Failed to delete product", http.StatusInternalServerError)
			return
		}

		if product == nil {
			response.SendJsonError(&w, "Product not found", http.StatusNotFound)
			return
		}

		err = h.ProductRepository.Delete(product.ID)

		if err != nil {
			response.SendJsonError(&w, "Failed to delete product", http.StatusInternalServerError)
			return
		}

		response.SendJsonSuccess(&w, "Product deleted successfully", product, http.StatusOK)
	}
}

func getUintId(w http.ResponseWriter, r *http.Request) (uint, error) {
	idString := r.PathValue("id")
	id, err := strconv.ParseUint(idString, 10, 32)

	if err != nil {
		response.SendJsonError(&w, "Invalid product id", http.StatusBadRequest)
		return 0, err
	}

	return uint(id), nil
}

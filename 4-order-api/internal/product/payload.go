package product

type CreateProductPayload struct {
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Images      []string `json:"images" validate:"required"`
}

type UpdateProductPayload struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
}

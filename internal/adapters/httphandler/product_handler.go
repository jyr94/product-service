package httphandler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/jyr94/product-service/internal/application"
	"github.com/jyr94/product-service/internal/domain"
)

type ProductHandler struct {
	service *application.ProductService
}

func NewProductHandler(service *application.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/products", h.createProduct).Methods(http.MethodPost)
	r.HandleFunc("/products", h.listProducts).Methods(http.MethodGet)
}

func (h *ProductHandler) createProduct(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.validateCreateProduct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	product := &domain.Product{
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
		Quantity:    req.Quantity,
	}

	if err := h.service.AddProduct(r.Context(), product); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "product created",
		"product_id": product.ID,
	})
}

func (h *ProductHandler) listProducts(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	sort := query.Get("sort")

	p := ParsePagination(r, 10)

	products, err := h.service.ListProducts(r.Context(), sort, p.Limit, p.Offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) validateCreateProduct(req CreateProductRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return errors.New("name is required")
	}
	if req.Price <= 0 {
		return errors.New("price must be greater than zero")
	}
	if req.Quantity < 0 {
		return errors.New("quantity must be zero or positive")
	}
	return nil
}

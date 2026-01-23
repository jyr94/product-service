package ports

import (
	"context"

	"github.com/jyr94/product-service/internal/domain"
)

type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	ListProducts(ctx context.Context, sort string, limit, offset int) ([]domain.Product, error)
}

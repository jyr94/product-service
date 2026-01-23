package application

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jyr94/product-service/internal/domain"
	"github.com/jyr94/product-service/internal/ports"
)

type ProductService struct {
	repo  ports.ProductRepository
	cache ports.Cache
}

func NewProductService(
	repo ports.ProductRepository,
	cache ports.Cache,
) *ProductService {
	return &ProductService{
		repo:  repo,
		cache: cache,
	}
}

func (s *ProductService) AddProduct(ctx context.Context, product *domain.Product) error {
	product.CreatedAt = time.Now()
	return s.repo.Create(ctx, product)
}

func (s *ProductService) ListProducts(ctx context.Context, sort string, limit int, offset int) ([]domain.Product, error) {

	cacheKey := buildProductListCacheKey(sort, limit, offset)

	if s.cache != nil {
		if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
			var products []domain.Product
			if err := json.Unmarshal(cached, &products); err == nil {
				return products, nil
			}
		}
	}

	products, err := s.repo.ListProducts(ctx, sort, limit, offset)
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		if bytes, err := json.Marshal(products); err == nil {
			_ = s.cache.Set(ctx, cacheKey, bytes, 60)
		}
	}

	return products, nil
}
func buildProductListCacheKey(sort string, limit, offset int) string {
	if sort == "" {
		sort = "latest"
	}
	return fmt.Sprintf("products:sort=%s:limit=%d:offset=%d", sort, limit, offset)
}

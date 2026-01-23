package application

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/jyr94/product-service/internal/domain"
	"github.com/jyr94/product-service/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductService_Create_Success(t *testing.T) {
	ctx := context.Background()

	repo := new(mocks.ProductRepository)
	cache := new(mocks.Cache)

	product := &domain.Product{
		Name:        "iPhone",
		Price:       20000000,
		Description: "phone",
		Quantity:    10,
	}

	repo.
		On("Create", ctx, product).
		Return(nil)

	service := NewProductService(repo, cache)

	err := service.AddProduct(ctx, product)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestProductService_Create_RepoError(t *testing.T) {
	ctx := context.Background()

	repo := new(mocks.ProductRepository)
	cache := new(mocks.Cache)

	product := &domain.Product{
		Name: "Macbook",
	}

	repo.
		On("Create", ctx, product).
		Return(errors.New("db error"))

	service := NewProductService(repo, cache)

	err := service.AddProduct(ctx, product)

	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestProductService_List_CacheHit(t *testing.T) {
	ctx := context.Background()

	expected := []domain.Product{
		{ID: 1, Name: "iPhone"},
	}

	cached, _ := json.Marshal(expected)

	repo := new(mocks.ProductRepository)
	cache := new(mocks.Cache)

	cache.
		On("Get", ctx, mock.Anything).
		Return(cached, nil)

	service := NewProductService(repo, cache)

	products, err := service.ListProducts(ctx, "latest", 10, 0)

	assert.NoError(t, err)
	assert.Len(t, products, 1)

	repo.AssertNotCalled(t, "ListProducts", mock.Anything)
}

func TestProductService_List_CacheMiss(t *testing.T) {
	ctx := context.Background()

	expected := []domain.Product{
		{ID: 2, Name: "Macbook"},
	}

	repo := new(mocks.ProductRepository)
	cache := new(mocks.Cache)

	cache.
		On("Get", ctx, mock.Anything).
		Return(nil, errors.New("cache miss"))

	repo.
		On("ListProducts", ctx, "latest", 10, 0).
		Return(expected, nil)

	cache.
		On("Set", ctx, mock.Anything, mock.Anything, mock.AnythingOfType("int")).
		Return(nil)

	service := NewProductService(repo, cache)

	products, err := service.ListProducts(ctx, "latest", 10, 0)

	assert.NoError(t, err)
	assert.Len(t, products, 1)

	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

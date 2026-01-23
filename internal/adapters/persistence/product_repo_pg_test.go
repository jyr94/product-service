package persistence

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jyr94/product-service/internal/domain"
)

func setUpRepoTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *ProductRepositoryPG) {
	db, mocks, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	repo := &ProductRepositoryPG{db: db}
	return db, mocks, repo
}

func TestProductRepository_List(t *testing.T) {
	db, mocks, repo := setUpRepoTest(t)
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"product_id",
		"product_name",
		"product_price",
		"product_description",
		"product_quantity",
		"created_at",
	}).AddRow(1, "Test Product", 9.99, "A product for testing", 100, time.Now())

	mocks.ExpectQuery(regexp.QuoteMeta(`
		SELECT
			product_id,
			product_name,
			product_price,
			product_description,
			product_quantity,
			created_at
		FROM products
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`)).
		WithArgs(10, 0).
		WillReturnRows(rows)

	products, err := repo.ListProducts(context.Background(), "latest", 10, 0)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(products) != 1 {
		t.Fatalf("expected 1 product, got %d", len(products))
	}

	if err := mocks.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}

}
func TestProductRepository_List_DBError(t *testing.T) {
	db, mock, repo := setUpRepoTest(t)
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT
			product_id,
			product_name,
			product_price,
			product_description,
			product_quantity,
			created_at
		FROM products
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`)).
		WithArgs(10, 0).
		WillReturnError(errors.New("db error"))

	_, err := repo.ListProducts(context.Background(), "latest", 10, 0)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
func TestProductRepository_add_success(t *testing.T) {
	db, mock, repo := setUpRepoTest(t)
	defer db.Close()

	product := domain.Product{
		ID:          1,
		Name:        "New Product",
		Price:       19.99,
		Description: "desc",
		Quantity:    10,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO products (product_name, product_price, product_description, product_quantity, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING product_id
	`)).
		WithArgs(

			product.Name,
			product.Price,
			product.Description,
			product.Quantity,
			product.CreatedAt,
		).
		WillReturnRows(
			sqlmock.NewRows([]string{"product_id"}).AddRow(1),
		)

	err := repo.Create(context.Background(), &product)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if product.ID != 1 {
		t.Fatalf("expected product ID to be set, got %d", product.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestProductRepository_Add_DBError(t *testing.T) {
	db, mock, repo := setUpRepoTest(t)
	defer db.Close()

	product := &domain.Product{
		Name:        "Macbook",
		Price:       30000000,
		Description: "laptop",
		Quantity:    5,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO products (product_name, product_price, product_description, product_quantity, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING product_id
	`)).
		WithArgs(
			product.Name,
			product.Price,
			product.Description,
			product.Quantity,
			product.CreatedAt,
		).
		WillReturnError(errors.New("db error"))

	err := repo.Create(context.Background(), product)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

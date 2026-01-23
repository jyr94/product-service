package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jyr94/product-service/internal/domain"
)

type ProductRepositoryPG struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepositoryPG {
	return &ProductRepositoryPG{db: db}
}

func (r *ProductRepositoryPG) Create(ctx context.Context, p *domain.Product) error {
	query := `
		INSERT INTO products (product_name, product_price, product_description, product_quantity, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING product_id
	`
	return r.db.QueryRowContext(
		ctx,
		query,
		p.Name,
		p.Price,
		p.Description,
		p.Quantity,
		p.CreatedAt,
	).Scan(&p.ID)
}

func (r *ProductRepositoryPG) ListProducts(ctx context.Context, sort string, limit, offset int) ([]domain.Product, error) {
	orderBy := "created_at DESC"

	switch sort {
	case "price_asc":
		orderBy = "product_price ASC"
	case "price_desc":
		orderBy = "product_price DESC"
	case "name_asc":
		orderBy = "product_name ASC"
	case "name_desc":
		orderBy = "product_name DESC"
	case "latest":
		orderBy = "created_at DESC"
	}

	query := fmt.Sprintf(`
		SELECT
			product_id,
			product_name,
			product_price,
			product_description,
			product_quantity,
			created_at
		FROM products
		ORDER BY %s
		LIMIT $1 OFFSET $2
	`, orderBy)

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Price,
			&p.Description,
			&p.Quantity,
			&p.CreatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

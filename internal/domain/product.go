package domain

import "time"

type Product struct {
	ID          int64
	Name        string
	Price       float64
	Description string
	Quantity    int
	CreatedAt   time.Time
}

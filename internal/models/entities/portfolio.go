package entities

import "time"

type Portfolio struct {
	ID        int64     `db:"id" json:"id"`
	UserID    int64     `db:"user_id" json:"user_id"`
	StockID   int64     `db:"stock_id" json:"stock_id"`
	Quantity  float64   `db:"quantity" json:"quantity"`
	Version   int       `db:"version" json:"version"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

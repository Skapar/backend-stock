package entities

import "time"

type Stock struct {
	ID        int64     `db:"id" json:"id"`
	Symbol    string    `db:"symbol" json:"symbol"`
	Name      string    `db:"name" json:"name"`
	Price     float64   `db:"price" json:"price"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

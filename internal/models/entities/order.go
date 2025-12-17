package entities

import "time"

type OrderType string

const (
	OrderBuy  OrderType = "BUY"
	OrderSell OrderType = "SELL"
)

type OrderStatus string

const (
	OrderPending   OrderStatus = "PENDING"
	OrderCompleted OrderStatus = "COMPLETED"
	OrderFailed    OrderStatus = "FAILED"
)

type Order struct {
	ID        int64       `db:"id" json:"id"`
	UserID    int64       `db:"user_id" json:"user_id"`
	StockID   int64       `db:"stock_id" json:"stock_id"`
	OrderType OrderType   `db:"order_type" json:"order_type"`
	Quantity  float64     `db:"quantity" json:"quantity"`
	Price     float64     `db:"price" json:"price"`
	Status    OrderStatus `db:"status" json:"status"`
	CreatedAt time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt time.Time   `db:"updated_at" json:"updated_at"`
}

package entities

import "time"

type HistoryAction string

const (
	ActionBuy           HistoryAction = "BUY"
	ActionSell          HistoryAction = "SELL"
	ActionBalanceUpdate HistoryAction = "BALANCE_UPDATE"
	ActionDeposit       HistoryAction = "DEPOSIT"
	ActionWithdraw      HistoryAction = "WITHDRAW"
)

type History struct {
	ID        int64         `db:"id" json:"id"`
	UserID    int64         `db:"user_id" json:"user_id"`
	OrderID   *int64        `db:"order_id" json:"order_id,omitempty"`
	StockID   *int64        `db:"stock_id" json:"stock_id,omitempty"`
	Action    HistoryAction `db:"action" json:"action"`
	Details   string        `db:"details" json:"details"`
	Amount    float64       `db:"amount" json:"amount"`
	CreatedAt time.Time     `db:"created_at" json:"created_at"`
}

package entities

import (
	"time"
)

type Subscription struct {
	ID             int64     `db:"id"`
	Name           string    `db:"name"`
	Price          int64     `db:"price"`
	DurationDays   int64     `db:"duration_days"`
	DailyLimitCups int64     `db:"daily_limit_cups"`
	TotalCups      int64     `db:"total_cups"`
	CreatedAt      time.Time `db:"created_at"`
}

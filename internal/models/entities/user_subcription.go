package entities

import "time"

type UserSubscription struct {
	ID             int64     `db:"id"`
	UserID         int64     `db:"user_id"`
	SubscriptionID int64     `db:"subscription_id"`
	StartDate      time.Time `db:"start_date"`
	EndDate        time.Time `db:"end_date"`
	TotalCups      int64     `db:"total_cups"`
	RemainingCups  int64     `db:"remaining_cups"`
}

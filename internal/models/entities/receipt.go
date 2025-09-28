package entities

import "time"

type Receipt struct {
	ID        int64      `db:"id"`
	UserID    int64      `db:"user_id"`
	FilePath  string     `db:"file_path"`
	Status    StatusType `db:"status"`
	CreatedAt time.Time  `db:"created_at"`
}

type StatusType string

const (
	StatusPending   StatusType = "pending"
	StatusApproved  StatusType = "approved"
	StatusRejected  StatusType = "rejected"
	StatusConfirmed StatusType = "confirmed"
)

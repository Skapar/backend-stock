package entities

import "time"

type Role string

const (
	RoleTrader Role = "TRADER"
	RoleAdmin  Role = "ADMIN"
)

type User struct {
	ID        int64     `db:"id"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Role      Role      `db:"role"`
	Balance   float64   `db:"balance"`
	CreatedAt time.Time `db:"created_at"`
}

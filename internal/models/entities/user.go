package entities

import "time"

type User struct {
	ID        int64     `db:"id"`
	TGID      int64     `db:"tg_id"`
	Username  string    `db:"username"`
	Name      string    `db:"name"`
	Phone     string    `db:"phone"`
	CreatedAt time.Time `db:"created_at"`
}

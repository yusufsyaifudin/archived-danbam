package user

import (
	"time"
)

type User struct {
	ID        string
	Enabled   bool
	CreatedAt time.Time
}

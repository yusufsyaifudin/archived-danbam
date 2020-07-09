package app

import (
	"time"
)

// App is table structure in database
type App struct {
	ID        string
	Name      string
	Enabled   bool
	CreatedAt time.Time
}

package models

import "time"

type Code struct {
	Value     string
	ExpiresAt time.Time
}

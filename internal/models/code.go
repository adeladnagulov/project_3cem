package models

import "time"

type Code struct {
	Value     int
	ExpiresAt time.Time
}

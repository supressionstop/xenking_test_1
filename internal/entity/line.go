package entity

import "time"

type Line struct {
	ID          int64
	Name        string
	Coefficient string
	SavedAt     time.Time
}

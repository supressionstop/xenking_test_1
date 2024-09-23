package entity

import "time"

type Line struct {
	Id          int64
	Name        string
	Coefficient string
	SavedAt     time.Time
}

package entity

import "time"

type Line struct {
	ID          int64
	Name        string
	Coefficient string
	SavedAt     time.Time
}

type LineMap map[Sport]Line

func LineMapFromSports(lines []Line) LineMap {
	m := make(LineMap, len(lines))
	for _, line := range lines {
		m[line.Name] = line
	}
	return m
}

type LinesDiff map[Sport]string

type Sport = string

const (
	Baseball Sport = "baseball"
	Football Sport = "football"
	Soccer   Sport = "soccer"
)

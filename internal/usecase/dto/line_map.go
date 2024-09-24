package dto

import (
	"github.com/supressionstop/xenking_test_1/internal/entity"
	"github.com/supressionstop/xenking_test_1/internal/usecase/enum"
)

type LineMap map[enum.Sport]entity.Line

func LineMapFromSports(lines []entity.Line) LineMap {
	m := make(LineMap, len(lines))
	for _, line := range lines {
		m[line.Name] = line
	}
	return m
}

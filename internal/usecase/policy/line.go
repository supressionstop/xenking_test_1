package policy

import (
	"xenking_test_1/internal/entity"
	"xenking_test_1/internal/usecase/dto"
)

func ProviderLineToEntity(line dto.ProviderLine) entity.Line {
	return entity.Line{
		Name:        line.Sport,
		Coefficient: line.Rate,
	}
}

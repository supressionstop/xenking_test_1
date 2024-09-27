package usecase

import (
	"github.com/stretchr/testify/assert"
	"github.com/supressionstop/xenking_test_1/internal/entity"
	"github.com/supressionstop/xenking_test_1/internal/usecase/dto"
	"testing"
)

func TestCalculateDiffUseCase(t *testing.T) {
	sut := NewCalculateDiffUseCase()

	tests := []struct {
		name    string
		prev    dto.LineMap
		curr    dto.LineMap
		got     dto.LinesDiff
		gotErr  error
		want    dto.LinesDiff
		wantErr error
	}{
		{
			name: "first single",
			prev: nil,
			curr: dto.LineMap{"soccer": entity.Line{
				Name:        "soccer",
				Coefficient: "1.5",
			}},
			want:    dto.LinesDiff{"soccer": "1.5"},
			wantErr: nil,
		},
		{
			name: "first multiple",
			prev: nil,
			curr: dto.LineMap{
				"soccer": entity.Line{
					Name:        "soccer",
					Coefficient: "1.5",
				},
				"football": entity.Line{
					Name:        "football",
					Coefficient: "0.555",
				},
			},
			want:    dto.LinesDiff{"soccer": "1.5", "football": "0.555"},
			wantErr: nil,
		},
		{
			name:    "diff single lesser",
			prev:    dto.LineMap{"soccer": entity.Line{Name: "soccer", Coefficient: "1.5"}},
			curr:    dto.LineMap{"soccer": entity.Line{Name: "soccer", Coefficient: "1.2"}},
			want:    dto.LinesDiff{"soccer": "-0.3"},
			wantErr: nil,
		},
		{
			name:    "diff single equal",
			prev:    dto.LineMap{"soccer": entity.Line{Name: "soccer", Coefficient: "1.5"}},
			curr:    dto.LineMap{"soccer": entity.Line{Name: "soccer", Coefficient: "1.5"}},
			want:    dto.LinesDiff{"soccer": "0"},
			wantErr: nil,
		},
		{
			name:    "diff single greater",
			prev:    dto.LineMap{"soccer": entity.Line{Name: "soccer", Coefficient: "1.5"}},
			curr:    dto.LineMap{"soccer": entity.Line{Name: "soccer", Coefficient: "2.5"}},
			want:    dto.LinesDiff{"soccer": "1"},
			wantErr: nil,
		},
		{
			name: "diff multiple",
			prev: dto.LineMap{
				"soccer": entity.Line{
					Name:        "soccer",
					Coefficient: "2.3",
				},
				"football": entity.Line{
					Name:        "football",
					Coefficient: "1.777",
				},
			},
			curr: dto.LineMap{
				"soccer": entity.Line{
					Name:        "soccer",
					Coefficient: "1.5",
				},
				"football": entity.Line{
					Name:        "football",
					Coefficient: "2.555",
				},
			},
			want:    dto.LinesDiff{"soccer": "-0.8", "football": "0.778"},
			wantErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// act
			got, gotErr := sut.Execute(tc.prev, tc.curr)

			// assert
			assert.Equal(t, tc.want, got)
			if tc.wantErr != nil {
				assert.EqualError(t, gotErr, tc.wantErr.Error())
			} else {
				assert.NoError(t, gotErr)
			}
		})
	}
}

// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package storage

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Baseball struct {
	ID        int32
	Sport     string
	Rate      pgtype.Numeric
	CreatedAt pgtype.Timestamp
}

type Football struct {
	ID        int32
	Sport     string
	Rate      pgtype.Numeric
	CreatedAt pgtype.Timestamp
}

type Soccer struct {
	ID        int32
	Sport     string
	Rate      pgtype.Numeric
	CreatedAt pgtype.Timestamp
}

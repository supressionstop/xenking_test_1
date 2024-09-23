package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
	"xenking_test_1/internal/entity"
	"xenking_test_1/internal/storage"
)

type Line struct {
	pg *storage.Postgres
}

func NewLine(pg *storage.Postgres) *Line {
	return &Line{pg: pg}
}

func (repo *Line) Save(ctx context.Context, line entity.Line) (entity.Line, error) {
	tx, err := repo.pg.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return entity.Line{}, fmt.Errorf("line.save begin tx: %w", err)
	}
	defer func() {
		err = repo.finishTransaction(ctx, err, tx)
	}()

	rate := pgtype.Numeric{}
	if err := rate.Scan(line.Coefficient); err != nil {
		return entity.Line{}, fmt.Errorf("string to numeric: %w", err)
	}
	savedLine, err := repo.pg.Queries.SaveLine(ctx, storage.SaveLineParams{
		Sport:     line.Name,
		Rate:      rate,
		CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
	})
	if err != nil {
		return entity.Line{}, fmt.Errorf("SaveLine: %w", err)
	}

	result := entity.Line{
		Id:      int64(savedLine.ID),
		Name:    savedLine.Sport,
		SavedAt: time.Time{},
	}
	newRate, err := savedLine.Rate.MarshalJSON()
	if err != nil {
		return entity.Line{}, fmt.Errorf("rate to string: %w", err)
	}
	result.Coefficient = string(newRate)

	return result, nil
}

func (repo *Line) finishTransaction(ctx context.Context, err error, tx pgx.Tx) error {
	if err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return fmt.Errorf("%s: %s", rollbackErr, err)
		}

		return err
	} else {
		if commitErr := tx.Commit(ctx); commitErr != nil {
			return fmt.Errorf("failed to commit tx: %s", err)
		}

		return nil
	}
}

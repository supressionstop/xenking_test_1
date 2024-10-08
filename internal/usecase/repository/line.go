package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	storage2 "github.com/supressionstop/xenking_test_1/internal/infrastructure/storage"
	"github.com/supressionstop/xenking_test_1/internal/usecase/entity"
)

type Line struct {
	pg *storage2.Postgres
}

func NewLine(pg *storage2.Postgres) *Line {
	return &Line{pg: pg}
}

func (repo *Line) Save(ctx context.Context, line entity.Line) (entity.Line, error) {
	tx, err := repo.pg.Pool.BeginTx(
		ctx,
		pgx.TxOptions{
			IsoLevel:       "",
			AccessMode:     "",
			DeferrableMode: "",
			BeginQuery:     "",
		},
	)
	if err != nil {
		return entity.Line{}, fmt.Errorf("line.save begin tx: %w", err)
	}

	defer func() {
		err = repo.finishTransaction(ctx, err, tx)
	}()

	sl, err := repo.saveLine(ctx, line)
	if err != nil {
		return entity.Line{}, fmt.Errorf("saving line: %w", err)
	}

	result, err := repo.savedLineToEntity(sl)
	if err != nil {
		return entity.Line{}, fmt.Errorf("hydration: %w", err)
	}

	return result, nil
}

func (repo *Line) GetRecentLines(ctx context.Context, sports []entity.Sport) ([]entity.Line, error) {
	tx, err := repo.pg.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("line.save begin tx: %w", err)
	}
	defer func() {
		err = repo.finishTransaction(ctx, err, tx)
	}()

	lines, err := repo.grabRecentLines(ctx, sports)
	if err != nil {
		return nil, fmt.Errorf("getting recent lines: %w", err)
	}

	result := make([]entity.Line, len(lines))
	for i, line := range lines {
		result[i], err = repo.savedLineToEntity(line)
		if err != nil {
			return nil, fmt.Errorf("hydration: %w", err)
		}
	}

	return result, nil
}

func (repo *Line) finishTransaction(ctx context.Context, err error, tx pgx.Tx) error {
	if err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return fmt.Errorf("failed to rollback: %w", rollbackErr)
		}

		return err
	} else {
		if commitErr := tx.Commit(ctx); commitErr != nil {
			return fmt.Errorf("failed to commit tx: %w", err)
		}

		return nil
	}
}

type savedLine struct {
	ID        int32
	Sport     string
	Rate      pgtype.Numeric
	CreatedAt pgtype.Timestamp
}

func (repo *Line) saveLine(ctx context.Context, line entity.Line) (savedLine, error) {
	var sl savedLine
	var errReturn error

	rate := pgtype.Numeric{}
	if err := rate.Scan(line.Coefficient); err != nil {
		return sl, fmt.Errorf("string to numeric: %w", err)
	}

	if line.Name == entity.Baseball {
		s, err := repo.pg.Queries.SaveBaseball(ctx, storage2.SaveBaseballParams{
			Sport:     line.Name,
			Rate:      rate,
			CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
		})
		sl = savedLine(s)
		errReturn = err
	} else if line.Name == entity.Football {
		s, err := repo.pg.Queries.SaveFootball(ctx, storage2.SaveFootballParams{
			Sport:     line.Name,
			Rate:      rate,
			CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
		})
		sl = savedLine(s)
		errReturn = err
	} else if line.Name == entity.Soccer {
		s, err := repo.pg.Queries.SaveSoccer(ctx, storage2.SaveSoccerParams{
			Sport:     line.Name,
			Rate:      rate,
			CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
		})
		sl = savedLine(s)
		errReturn = err
	} else {
		return sl, fmt.Errorf("unknown line")
	}

	return sl, errReturn
}

func (repo *Line) savedLineToEntity(sl savedLine) (entity.Line, error) {
	result := entity.Line{
		ID:      int64(sl.ID),
		Name:    sl.Sport,
		SavedAt: time.Time{},
	}
	newRate, err := sl.Rate.MarshalJSON()
	if err != nil {
		return entity.Line{}, fmt.Errorf("rate to string: %w", err)
	}
	result.Coefficient = string(newRate)

	return result, nil
}

func (repo *Line) grabRecentLines(ctx context.Context, sports []entity.Sport) ([]savedLine, error) {
	var errReturn error

	result := make([]savedLine, len(sports))

	for idx, sport := range sports {
		if sport == entity.Baseball {
			s, err := repo.pg.Queries.GetRecentBaseball(ctx)
			result[idx] = savedLine(s)
			errReturn = err
		} else if sport == entity.Football {
			s, err := repo.pg.Queries.GetRecentFootball(ctx)
			result[idx] = savedLine(s)
			errReturn = err
		} else if sport == entity.Soccer {
			s, err := repo.pg.Queries.GetRecentSoccer(ctx)
			result[idx] = savedLine(s)
			errReturn = err
		} else {
			return nil, fmt.Errorf("unknown sport")
		}
	}

	return result, errReturn
}

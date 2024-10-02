package storage

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"

	_ "github.com/jackc/pgx/v5/stdlib" // sql driver init
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

const (
	defaultConnAttempts = 10
	defaultConnTimeout  = time.Second
)

type Postgres struct {
	Pool         *pgxpool.Pool
	Queries      *Queries
	logger       *slog.Logger
	connAttempts int
	connTimeout  time.Duration
}

// NewPostgres requires to close Pool after usage.
// Uses pgx.Pool to handle concurrent operations.
func NewPostgres(ctx context.Context, connString string, logger *slog.Logger) (*Postgres, error) {
	pg := &Postgres{
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout,
		logger:       logger,
	}

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	pg.Pool, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create Pool: %w", err)
	}

	if err := pg.tryPing(ctx); err != nil {
		return nil, fmt.Errorf("connect to db: %w", err)
	}

	pg.Queries = New(pg.Pool)

	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return nil, fmt.Errorf("set migrations dialect: %w", err)
	}

	return pg, nil
}

func (pg *Postgres) ClosePool() {
	pg.Pool.Close()
}

// Up applies migrations to DB.
func (pg *Postgres) Up(connString string) error {
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return fmt.Errorf("up migrations: %w", err)
	}
	defer db.Close()

	return goose.Up(db, "migrations")
}

func (pg *Postgres) tryPing(ctx context.Context) error {
	pingCh := make(chan error)
	go func(resultCh chan error) {
		for pg.connAttempts > 0 {
			err := pg.Pool.Ping(ctx)
			if err == nil {
				pg.logger.Info("db ping succeeded")
				resultCh <- nil
				break
			}
			pg.logger.Error(
				"ping db failed, retrying",
				slog.Any("err", err),
				slog.Int("attempts_left", pg.connAttempts),
			)
			time.Sleep(pg.connTimeout)
			pg.connAttempts--
		}
		resultCh <- fmt.Errorf("ping failed")
	}(pingCh)

	select {
	case err := <-pingCh:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

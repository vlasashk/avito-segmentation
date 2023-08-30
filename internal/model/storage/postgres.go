package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

type PostgresDB struct {
	DB *pgxpool.Pool
}

func New(ctx context.Context) (*PostgresDB, error) {
	dbPool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v\n", err)
	}
	pgInstance := &PostgresDB{dbPool}
	if err = pgInstance.Migrate(ctx); err != nil {
		return nil, err
	}
	return pgInstance, nil
}

func (pg *PostgresDB) DropDB(ctx context.Context) error {
	err := pg.DB.AcquireFunc(context.Background(), func(conn *pgxpool.Conn) error {
		if err := pg.Ping(ctx); err != nil {
			return fmt.Errorf("failed to ping db: %v\n", err)
		}
		query, err := os.ReadFile(os.Getenv("DROP_PATH"))
		if err != nil {
			return fmt.Errorf("failed to read drop.sql file: %v\n", err)
		}
		if _, err = conn.Exec(ctx, string(query)); err != nil {
			return fmt.Errorf("failed to drop tables: %v\n", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("unable to acquire a database connection: %v\n", err)
	}
	return nil
}

func (pg *PostgresDB) Migrate(ctx context.Context) error {
	err := pg.DB.AcquireFunc(context.Background(), func(conn *pgxpool.Conn) error {
		if err := pg.Ping(ctx); err != nil {
			return fmt.Errorf("failed to ping db: %v\n", err)
		}
		query, err := os.ReadFile(os.Getenv("MIGRATE_PATH"))
		if err != nil {
			return fmt.Errorf("failed to read up.sql file: %v\n", err)
		}
		if _, err = conn.Exec(ctx, string(query)); err != nil {
			return fmt.Errorf("failed to init tables: %v\n", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("unable to acquire a database connection: %v\n", err)
	}
	return nil
}

func (pg *PostgresDB) Ping(ctx context.Context) error {
	return pg.DB.Ping(ctx)
}

func (pg *PostgresDB) Close() {
	pg.DB.Close()
}

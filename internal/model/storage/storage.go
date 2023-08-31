package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type User struct {
	Id  uint64 `json:"id,omitempty"`
	UID uint64 `json:"user_id" validate:"required"`
}

type Segment struct {
	Id   uint64 `json:"id,omitempty"`
	Slug string `json:"slug" validate:"required"`
}

type UserSegment struct {
	UserID    uint64    `json:"user_id"`
	SegmentID uint64    `json:"segment_id"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
}

func (pg *PostgresDB) CascadeDeleteSegment() {

}
func (pg *PostgresDB) DeleteUserFromSegment() {

}
func (pg *PostgresDB) AddUserToSegment() {

}
func (pg *PostgresDB) GetUserSegmentsInfo() {

}
func (pg *PostgresDB) GetSegmentUsersInfo() {

}
func (pg *PostgresDB) AddUser(ctx context.Context, user User) (uint64, error) {
	var id uint64
	err := pg.DB.AcquireFunc(context.Background(), func(conn *pgxpool.Conn) error {
		if err := pg.Ping(ctx); err != nil {
			return fmt.Errorf("failed to ping db: %v\n", err)
		}
		query := `insert into users (user_id) values ($1) returning "id"`
		if err := conn.QueryRow(ctx, query, user.UID).Scan(&id); err != nil {
			return fmt.Errorf("failed to insert data: %v\n", err)
		}
		return nil
	})
	if err != nil {
		return id, err
	}
	return id, nil
}

func (pg *PostgresDB) AddSegment(ctx context.Context, segment Segment) (uint64, error) {
	var id uint64
	err := pg.DB.AcquireFunc(context.Background(), func(conn *pgxpool.Conn) error {
		if err := pg.Ping(ctx); err != nil {
			return fmt.Errorf("failed to ping db: %v\n", err)
		}
		query := `insert into segments (slug) values ($1) returning "id"`
		if err := conn.QueryRow(ctx, query, segment.Slug).Scan(&id); err != nil {
			return fmt.Errorf("failed to insert data: %v\n", err)
		}
		return nil
	})
	if err != nil {
		return id, err
	}
	return id, nil
}

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
	UserID      uint64    `json:"user_id" validate:"required"`
	SegmentSlug string    `json:"segment_slug" validate:"required"`
	CreatedAt   time.Time `json:"created_at"`
	DeletedAt   time.Time `json:"deleted_at,omitempty"`
}

func (pg *PostgresDB) CascadeDeleteSegment() {

}
func (pg *PostgresDB) DeleteUserFromSegment() {

}
func (pg *PostgresDB) GetSegmentUsersInfo() {

}

func (pg *PostgresDB) GetUserSegmentsInfo(ctx context.Context, user User) ([]string, error) {
	var id uint64
	var res []string
	err := pg.DB.AcquireFunc(context.Background(), func(conn *pgxpool.Conn) error {
		if err := pg.Ping(ctx); err != nil {
			return fmt.Errorf("failed to ping db: %v\n", err)
		}
		queryCheckUser := `select id from users where user_id = $1`
		query := `select slug from user_segments us join segments s on s.id = us.segment_id where user_id = $1`
		if err := conn.QueryRow(ctx, queryCheckUser, user.UID).Scan(&id); err != nil {
			return fmt.Errorf("user %v doesn't exist: %v\n", user.UID, err)
		}
		if rows, err := conn.Query(ctx, query, id); err != nil {
			return fmt.Errorf("failed to get data: %v\n", err)
		} else {
			defer rows.Close()
			for rows.Next() {
				var segment string
				if err = rows.Scan(&segment); err != nil {
					return fmt.Errorf("failed to scan segment: %v\n", err)
				}
				res = append(res, segment)
			}
			if err = rows.Err(); err != nil {
				return fmt.Errorf("error occurred while reading: %v\n", err)
			}
		}
		return nil
	})
	if err != nil {
		return res, err
	}
	return res, nil
}

func (pg *PostgresDB) AddUserToSegment(ctx context.Context, userSegment UserSegment) error {
	var userID, segmentID uint64
	err := pg.DB.AcquireFunc(context.Background(), func(conn *pgxpool.Conn) error {
		if err := pg.Ping(ctx); err != nil {
			return fmt.Errorf("failed to ping db: %v\n", err)
		}
		queryCheckUser := `select id from users where user_id = $1`
		queryCheckSegment := `select id from segments where slug = $1`
		queryInsert := `insert into user_segments (user_id, segment_id) values ($1, $2)`
		if err := conn.QueryRow(ctx, queryCheckUser, userSegment.UserID).Scan(&userID); err != nil {
			return fmt.Errorf("user %v doesn't exist: %v\n", userSegment.UserID, err)
		}
		if err := conn.QueryRow(ctx, queryCheckSegment, userSegment.SegmentSlug).Scan(&segmentID); err != nil {
			return fmt.Errorf("segment %v doesn't exist: %v\n", userSegment.SegmentSlug, err)
		}
		if _, err := conn.Exec(ctx, queryInsert, userID, segmentID); err != nil {
			return fmt.Errorf("failed to insert data: %v\n", err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
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

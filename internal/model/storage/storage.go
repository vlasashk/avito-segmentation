package storage

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vlasashk/avito-segmentation/internal/model/logger"
	"log/slog"
	"os"
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

type UserSegments struct {
	UserID      uint64   `json:"user_id" validate:"required"`
	SegmentSlug []string `json:"segment_slug" validate:"required"`
}

type CsvReport struct {
	Year  uint       `json:"year" validate:"required"`
	Month time.Month `json:"month" validate:"required"`
}

func (pg *PostgresDB) CsvHistoryReport(ctx context.Context, csvDates CsvReport, log *slog.Logger) error {
	err := pg.DB.AcquireFunc(context.Background(), func(conn *pgxpool.Conn) error {
		if err := pg.Ping(ctx); err != nil {
			log.Error("failed to ping db", logger.Err(err))
			return fmt.Errorf("failed to ping db")
		}
		startDate := time.Date(int(csvDates.Year), csvDates.Month, 0, 0, 0, 0, 0, time.Local)
		endDate := startDate.AddDate(0, 1, 0)
		query := `SELECT u.user_id, s.slug AS segment, 'added' AS status, created_at AS segment_date
					FROM user_segments us
					JOIN segments s on s.id = us.segment_id
					JOIN users u on u.id = us.user_id
					WHERE (created_at BETWEEN $1 AND $2)
					UNION
					SELECT u.user_id, s.slug AS segment, 'removed' AS status, deleted_at AS segment_date
					FROM user_segments us
					JOIN segments s on s.id = us.segment_id
					JOIN users u on u.id = us.user_id
					WHERE (deleted_at BETWEEN $1 AND $2)
					ORDER BY user_id, segment, status;`
		var file *os.File
		reportPath := fmt.Sprintf("%s%s_%d_%d.csv", os.Getenv("CSV_PATH"), "report", csvDates.Year, csvDates.Month)
		if err := os.Mkdir(os.Getenv("CSV_PATH"), 0755); err != nil && !os.IsExist(err) {
			log.Error("failed to create directory", logger.Err(err))
			return fmt.Errorf("failed to create directory")
		}
		if fileTemp, err := os.Create(reportPath); err != nil {
			log.Error("failed to create file", logger.Err(err))
			return fmt.Errorf("failed to create file")
		} else {
			file = fileTemp
		}
		defer func(file *os.File) {
			_ = file.Close()
		}(file)

		if rows, err := conn.Query(ctx, query, startDate, endDate); err != nil {
			log.Error("failed to execute query", logger.Err(err))
			return fmt.Errorf("failed to execute query")
		} else {
			defer rows.Close()
			if err = writeCsv(file, rows, log); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func writeCsv(file *os.File, rows pgx.Rows, log *slog.Logger) error {
	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()
	for rows.Next() {
		var userId, segmentId, status string
		var segmentDate time.Time
		if err := rows.Scan(&userId, &segmentId, &status, &segmentDate); err != nil {
			log.Error("failed to scan data", logger.Err(err))
			return fmt.Errorf("failed to scan data")
		}
		if err := writer.Write([]string{userId, segmentId, status, segmentDate.String()}); err != nil {
			log.Error("failed to write into csv:", logger.Err(err))
			return fmt.Errorf("failed to write into csv")
		}
	}
	if err := rows.Err(); err != nil {
		log.Error("error occurred while reading", logger.Err(err))
		return fmt.Errorf("error occurred while reading")
	}
	return nil
}

func (pg *PostgresDB) GetSegmentUsersInfo(ctx context.Context, segment Segment, log *slog.Logger) ([]uint64, error) {
	var id uint64
	var res []uint64
	err := pg.DB.AcquireFunc(context.Background(), func(conn *pgxpool.Conn) error {
		if err := pg.Ping(ctx); err != nil {
			log.Error("failed to ping db", logger.Err(err))
			return fmt.Errorf("failed to ping db")
		}
		queryCheckUser := `select id from segments where slug = $1`
		query := `select u.user_id from user_segments us join users u on u.id = us.user_id where segment_id = $1 and deleted_at is null`
		if err := conn.QueryRow(ctx, queryCheckUser, segment.Slug).Scan(&id); err != nil {
			log.Error(fmt.Sprintf("segment '%v' doesn't exist", segment.Slug), logger.Err(err))
			return fmt.Errorf("segment '%v' doesn't exist", segment.Slug)
		}
		if rows, err := conn.Query(ctx, query, id); err != nil {
			log.Error("failed to get data", logger.Err(err))
			return fmt.Errorf("failed to get data")
		} else {
			defer rows.Close()
			for rows.Next() {
				var user uint64
				if err = rows.Scan(&user); err != nil {
					log.Error("failed to scan user", logger.Err(err))
					return fmt.Errorf("failed to scan user")
				}
				res = append(res, user)
			}
			if err = rows.Err(); err != nil {
				log.Error("error occurred while reading", logger.Err(err))
				return fmt.Errorf("error occurred while reading")
			}
		}
		return nil
	})
	if err != nil {
		return res, err
	}
	return res, nil
}

func (pg *PostgresDB) DeleteUserFromSegments(ctx context.Context, userSegment UserSegments, log *slog.Logger) error {
	var id uint64
	err := pg.DB.AcquireFunc(context.Background(), func(conn *pgxpool.Conn) error {
		if err := pg.Ping(ctx); err != nil {
			log.Error("failed to ping db", logger.Err(err))
			return fmt.Errorf("failed to ping db")
		}
		queryCheckUser := `select id from users where user_id = $1`
		query := `update user_segments us
				  set deleted_at = NOW()
				  where segment_id in (select id from segments where slug = $1)
  					and us.user_id = $2
  					and deleted_at is null;`
		if err := conn.QueryRow(ctx, queryCheckUser, userSegment.UserID).Scan(&id); err != nil {
			log.Error(fmt.Sprintf("user '%v' doesn't exist", userSegment.UserID), logger.Err(err))
			return fmt.Errorf("user '%v' doesn't exist", userSegment.UserID)
		}
		if tx, err := conn.Begin(ctx); err != nil {
			log.Error("failed to begin transaction", logger.Err(err))
			return fmt.Errorf("failed to begin transaction")
		} else {
			defer func(tx pgx.Tx, ctx context.Context) {
				_ = tx.Rollback(ctx)
			}(tx, context.Background())
			for i := 0; i < len(userSegment.SegmentSlug); i++ {
				if _, err = conn.Exec(ctx, query, userSegment.SegmentSlug[i], id); err != nil {
					log.Error("failed to delete data", logger.Err(err))
					return fmt.Errorf("failed to delete data")
				}
			}
			err = tx.Commit(context.Background())
			if err != nil {
				log.Error("failed to commit transaction", logger.Err(err))
				return fmt.Errorf("failed to commit transaction")
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (pg *PostgresDB) AddUserToSegments(ctx context.Context, userSegment UserSegments, log *slog.Logger) error {
	var userID, segmentID uint64
	segmentSlice := make([]uint64, 0, len(userSegment.SegmentSlug))
	err := pg.DB.AcquireFunc(context.Background(), func(conn *pgxpool.Conn) error {
		if err := pg.Ping(ctx); err != nil {
			log.Error("failed to ping db", logger.Err(err))
			return fmt.Errorf("failed to ping db")
		}
		queryCheckUser := `select id from users where user_id = $1`
		queryCheckSegment := `select id from segments where slug = $1`
		queryInsert := `insert into user_segments (user_id, segment_id)
						values ($1, $2)
						on conflict (user_id, segment_id) do update set deleted_at = NULL, created_at = NOW()
						where user_segments.segment_id in (select id from segments where slug = $3)
						  and user_segments.user_id = $1
						  and user_segments.deleted_at is not null;`
		if err := conn.QueryRow(ctx, queryCheckUser, userSegment.UserID).Scan(&userID); err != nil {
			log.Error(fmt.Sprintf("user '%v' doesn't exist", userSegment.UserID), logger.Err(err))
			return fmt.Errorf("user '%v' doesn't exist", userSegment.UserID)
		}
		for i := 0; i < len(userSegment.SegmentSlug); i++ {
			if err := conn.QueryRow(ctx, queryCheckSegment, userSegment.SegmentSlug[i]).Scan(&segmentID); err != nil {
				log.Error(fmt.Sprintf("segment '%v' doesn't exist", userSegment.SegmentSlug), logger.Err(err))
				return fmt.Errorf("segment '%v' doesn't exist", userSegment.SegmentSlug[i])
			}
			segmentSlice = append(segmentSlice, segmentID)
		}
		if tx, err := conn.Begin(ctx); err != nil {
			log.Error("failed to begin transaction", logger.Err(err))
			return fmt.Errorf("failed to begin transaction")
		} else {
			defer func(tx pgx.Tx, ctx context.Context) {
				_ = tx.Rollback(ctx)
			}(tx, context.Background())
			for i := 0; i < len(segmentSlice); i++ {
				var res pgconn.CommandTag
				if res, err = conn.Exec(ctx, queryInsert, userID, segmentSlice[i], userSegment.SegmentSlug[i]); err != nil {
					log.Error("failed to insert data or data already exists", logger.Err(err))
					return fmt.Errorf("failed to insert data or data already exists")
				} else {
					n := res.RowsAffected()
					if n < 1 {
						log.Error("failed to execute query", logger.Err(fmt.Errorf("failed to update data or data is already present in DB")))
						return fmt.Errorf("failed to update data or data already exists")
					}
				}
			}
			err = tx.Commit(context.Background())
			if err != nil {
				log.Error("failed to commit transaction", logger.Err(err))
				return fmt.Errorf("failed to commit transaction")
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (pg *PostgresDB) GetUserSegmentsInfo(ctx context.Context, user User, log *slog.Logger) ([]string, error) {
	var id uint64
	var res []string
	err := pg.DB.AcquireFunc(context.Background(), func(conn *pgxpool.Conn) error {
		if err := pg.Ping(ctx); err != nil {
			log.Error("failed to ping db", logger.Err(err))
			return fmt.Errorf("failed to ping db")
		}
		queryCheckUser := `select id from users where user_id = $1`
		query := `select slug from user_segments us join segments s on s.id = us.segment_id where user_id = $1 and deleted_at is null `
		if err := conn.QueryRow(ctx, queryCheckUser, user.UID).Scan(&id); err != nil {
			log.Error(fmt.Sprintf("user '%v' doesn't exist", user.UID), logger.Err(err))
			return fmt.Errorf("user '%v' doesn't exist", user.UID)
		}
		if rows, err := conn.Query(ctx, query, id); err != nil {
			log.Error("failed to get data", logger.Err(err))
			return fmt.Errorf("failed to get data")
		} else {
			defer rows.Close()
			for rows.Next() {
				var segment string
				if err = rows.Scan(&segment); err != nil {
					log.Error("failed to scan segment", logger.Err(err))
					return fmt.Errorf("failed to scan segment")
				}
				res = append(res, segment)
			}
			if err = rows.Err(); err != nil {
				log.Error("error occurred while reading", logger.Err(err))
				return fmt.Errorf("error occurred while reading")
			}
		}
		return nil
	})
	if err != nil {
		return res, err
	}
	return res, nil
}

func (pg *PostgresDB) AddUser(ctx context.Context, user User, log *slog.Logger) (uint64, error) {
	var id uint64
	err := pg.DB.AcquireFunc(context.Background(), func(conn *pgxpool.Conn) error {
		if err := pg.Ping(ctx); err != nil {
			log.Error("failed to ping db", logger.Err(err))
			return fmt.Errorf("failed to ping db")
		}
		query := `insert into users (user_id) values ($1) returning "id"`
		if err := conn.QueryRow(ctx, query, user.UID).Scan(&id); err != nil {
			log.Error("failed to insert data or user already exists", logger.Err(err))
			return fmt.Errorf("failed to insert data or user already exists")
		}
		return nil
	})
	if err != nil {
		return id, err
	}
	return id, nil
}

func (pg *PostgresDB) AddSegment(ctx context.Context, segment Segment, log *slog.Logger) (uint64, error) {
	var id uint64
	err := pg.DB.AcquireFunc(context.Background(), func(conn *pgxpool.Conn) error {
		if err := pg.Ping(ctx); err != nil {
			log.Error("failed to ping db", logger.Err(err))
			return fmt.Errorf("failed to ping db")
		}
		query := `insert into segments (slug) values ($1) returning "id"`
		if err := conn.QueryRow(ctx, query, segment.Slug).Scan(&id); err != nil {
			log.Error("failed to insert data or segment already exists", logger.Err(err))
			return fmt.Errorf("failed to insert data or segment already exists")
		}
		return nil
	})
	if err != nil {
		return id, err
	}
	return id, nil
}

func (pg *PostgresDB) CascadeDeleteSegment(ctx context.Context, segment Segment, log *slog.Logger) error {
	err := pg.DB.AcquireFunc(context.Background(), func(conn *pgxpool.Conn) error {
		if err := pg.Ping(ctx); err != nil {
			log.Error("failed to ping db", logger.Err(err))
			return fmt.Errorf("failed to ping db")
		}
		query := `delete from segments where slug = $1`
		if res, err := conn.Exec(ctx, query, segment.Slug); err != nil {
			log.Error("failed to delete segment", logger.Err(err))
			return fmt.Errorf("failed to delete segment")
		} else {
			n := res.RowsAffected()
			if n < 1 {
				log.Error("failed to execute query", logger.Err(fmt.Errorf("segment %v is not present in database\n", segment.Slug)))
				return fmt.Errorf("segment %v is not present in database\n", segment.Slug)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

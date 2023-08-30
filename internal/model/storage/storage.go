package storage

import (
	"time"
)

type Storage interface {
	Select()
	Get()
	Post()
	Delete()
}

type User struct {
	Id  uint64
	UID uint64
}

type Segment struct {
	Id   uint64
	Slug string
}

type UserSegment struct {
	UserID    uint64
	SegmentID uint64
	CreatedAt time.Time
	DeletedAt time.Time
}

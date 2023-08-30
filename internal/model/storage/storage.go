package storage

import (
	"time"
)

type Storage interface {
	CascadeDeleteSegment()
	DeleteUserFromSegment()
	AddUserToSegment()
	GetUserSegmentsInfo()
	GetSegmentUsersInfo()
	AddUser(uid uint64)
	AddSegment(slug string)
}

type User struct {
	Id  uint64 `json:"id"`
	UID uint64 `json:"user_id"`
}

type Segment struct {
	Id   uint64 `json:"id"`
	Slug string `json:"slug"`
}

type UserSegment struct {
	UserID    uint64    `json:"user_id"`
	SegmentID uint64    `json:"segment_id"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

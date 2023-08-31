package api

import (
	"context"
	"github.com/vlasashk/avito-segmentation/internal/model/storage"
	"log/slog"
)

type Storage interface {
	CascadeDeleteSegment()
	DeleteUserFromSegment()
	GetUserSegmentsInfo(context.Context, storage.User) ([]string, error)
	GetSegmentUsersInfo()
	AddUserToSegment(context.Context, storage.UserSegment) error
	AddUser(context.Context, storage.User) (uint64, error)
	AddSegment(context.Context, storage.Segment) (uint64, error)
}

type ServerAPI struct {
	ListenAddr string
	Store      Storage
	Log        *slog.Logger
}

type UserRequest struct {
	storage.User
}

type SegmentRequest struct {
	storage.Segment
}

type UserSegmentRequest struct {
	storage.UserSegment
}

type UserResponse struct {
	ResponseStatus
	storage.User
}

type SegmentResponse struct {
	ResponseStatus
	storage.Segment
}

type UserSegmentResponse struct {
	ResponseStatus
	storage.UserSegment
}

type GetSegmentsResponse struct {
	ResponseStatus
	UserID      uint64   `json:"user_id" validate:"required"`
	SegmentSlug []string `json:"user_segments" validate:"required"`
}

type GetUsersResponse struct {
	ResponseStatus
	SegmentSlug string   `json:"user_segment" validate:"required"`
	UserIDs     []uint64 `json:"user_ids" validate:"required"`
}

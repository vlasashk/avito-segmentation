package api

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/vlasashk/avito-segmentation/internal/model/logger"
	"github.com/vlasashk/avito-segmentation/internal/model/storage"
	"log/slog"
	"net/http"
)

type StorageHandler interface {
	CascadeDeleteSegment()
	DeleteUserFromSegment()
	AddUserToSegment()
	GetUserSegmentsInfo()
	GetSegmentUsersInfo()
	AddUser(uid uint64)
	AddSegment(slug string)
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

func (s *ServerAPI) HandleAddUser(w http.ResponseWriter, r *http.Request) {
	newUser := &UserRequest{}
	log := s.Log.With(
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	if err := render.DecodeJSON(r.Body, &newUser); err != nil {
		log.Error("failed to decode request body", logger.Err(err))
		render.JSON(w, r, Error("failed to decode request body"))
		return
	}
	log.Info("request body decoded", slog.Any("request", *newUser))
	if err := validator.New().Struct(newUser); err != nil {
		log.Error("wrong body structure", logger.Err(err))
		render.JSON(w, r, Error("wrong body structure"))
		return
	}
	id, err := s.Store.AddUser(context.Background(), newUser.User)
	if err != nil {
		log.Error("failed to execute query", logger.Err(err))
		render.JSON(w, r, Error("failed to execute query"))
		return
	}
	response := UserResponse{
		ResponseStatus: OK(),
		User:           newUser.User,
	}
	response.User.Id = id
	log.Info("query successfully executed", slog.Any("request", response))
	render.JSON(w, r, response)
	return
}

func (s *ServerAPI) HandleAddSegment(w http.ResponseWriter, r *http.Request) {
	newSegment := &SegmentRequest{}
	log := s.Log.With(
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	if err := render.DecodeJSON(r.Body, &newSegment); err != nil {
		log.Error("failed to decode request body", logger.Err(err))
		render.JSON(w, r, Error("failed to decode request body"))
		return
	}
	log.Info("request body decoded", slog.Any("request", *newSegment))
	if err := validator.New().Struct(newSegment); err != nil {
		log.Error("wrong body structure", logger.Err(err))
		render.JSON(w, r, Error("wrong body structure"))
		return
	}
	id, err := s.Store.AddSegment(context.Background(), newSegment.Segment)
	if err != nil {
		log.Error("failed to execute query", logger.Err(err))
		render.JSON(w, r, Error("failed to execute query"))
		return
	}
	response := SegmentResponse{
		ResponseStatus: OK(),
		Segment:        newSegment.Segment,
	}
	response.Segment.Id = id
	log.Info("query successfully executed", slog.Any("request", response))
	render.JSON(w, r, response)
	return
}

func (s *ServerAPI) HandleGetSegmentUsersInfo(w http.ResponseWriter, r *http.Request) {
	return
}

func (s *ServerAPI) HandleGetUserSegmentsInfo(w http.ResponseWriter, r *http.Request) {
	return
}

func (s *ServerAPI) HandleAddUserToSegment(w http.ResponseWriter, r *http.Request) {
	return
}

func (s *ServerAPI) HandleDeleteUserFromSegment(w http.ResponseWriter, r *http.Request) {
	return
}

func (s *ServerAPI) HandleCascadeDeleteSegment(w http.ResponseWriter, r *http.Request) {
	return
}

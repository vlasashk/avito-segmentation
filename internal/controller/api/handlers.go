package api

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/vlasashk/avito-segmentation/internal/model/logger"
	"log/slog"
	"net/http"
)

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

func (s *ServerAPI) HandleAddUserToSegment(w http.ResponseWriter, r *http.Request) {
	newUserSegment := &UserSegmentRequest{}
	log := s.Log.With(
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	if err := render.DecodeJSON(r.Body, &newUserSegment); err != nil {
		log.Error("failed to decode request body", logger.Err(err))
		render.JSON(w, r, Error("failed to decode request body"))
		return
	}
	log.Info("request body decoded", slog.Any("request", *newUserSegment))
	if err := validator.New().Struct(newUserSegment); err != nil {
		log.Error("wrong body structure", logger.Err(err))
		render.JSON(w, r, Error("wrong body structure"))
		return
	}
	if len(newUserSegment.UserSegments.SegmentSlug) == 0 {
		log.Error("empty segment array")
		render.JSON(w, r, Error("empty segment array"))
		return
	}
	err := s.Store.AddUserToSegments(context.Background(), newUserSegment.UserSegments)
	if err != nil {
		log.Error("failed to execute query", logger.Err(err))
		render.JSON(w, r, Error("failed to execute query"))
		return
	}
	response := UserSegmentResponse{
		ResponseStatus: OK(),
		UserSegments:   newUserSegment.UserSegments,
	}
	log.Info("query successfully executed", slog.Any("request", response))
	render.JSON(w, r, response)
	return
}

func (s *ServerAPI) HandleGetUserSegmentsInfo(w http.ResponseWriter, r *http.Request) {
	user := &UserRequest{}
	log := s.Log.With(
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	if err := render.DecodeJSON(r.Body, &user); err != nil {
		log.Error("failed to decode request body", logger.Err(err))
		render.JSON(w, r, Error("failed to decode request body"))
		return
	}
	log.Info("request body decoded", slog.Any("request", *user))
	if err := validator.New().Struct(user); err != nil {
		log.Error("wrong body structure", logger.Err(err))
		render.JSON(w, r, Error("wrong body structure"))
		return
	}
	segments, err := s.Store.GetUserSegmentsInfo(context.Background(), user.User)
	if err != nil {
		log.Error("failed to execute query", logger.Err(err))
		render.JSON(w, r, Error("failed to execute query"))
		return
	}
	response := GetSegmentsResponse{
		ResponseStatus: OK(),
		UserID:         user.UID,
		SegmentSlug:    segments,
	}
	log.Info("query successfully executed", slog.Any("request", response))
	render.JSON(w, r, response)
	return
}

func (s *ServerAPI) HandleDeleteUserFromSegment(w http.ResponseWriter, r *http.Request) {
	newUserSegment := &UserSegmentRequest{}
	log := s.Log.With(
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	if err := render.DecodeJSON(r.Body, &newUserSegment); err != nil {
		log.Error("failed to decode request body", logger.Err(err))
		render.JSON(w, r, Error("failed to decode request body"))
		return
	}
	log.Info("request body decoded", slog.Any("request", *newUserSegment))
	if err := validator.New().Struct(newUserSegment); err != nil {
		log.Error("wrong body structure", logger.Err(err))
		render.JSON(w, r, Error("wrong body structure"))
		return
	}
	if len(newUserSegment.UserSegments.SegmentSlug) == 0 {
		log.Error("empty segment array")
		render.JSON(w, r, Error("empty segment array"))
		return
	}
	err := s.Store.DeleteUserFromSegments(context.Background(), newUserSegment.UserSegments)
	if err != nil {
		log.Error("failed to execute query", logger.Err(err))
		render.JSON(w, r, Error("failed to execute query"))
		return
	}
	response := UserSegmentResponse{
		ResponseStatus: OK(),
		UserSegments:   newUserSegment.UserSegments,
	}
	log.Info("query successfully executed", slog.Any("request", response))
	render.JSON(w, r, response)
	return
}

func (s *ServerAPI) HandleCascadeDeleteSegment(w http.ResponseWriter, r *http.Request) {
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
	err := s.Store.CascadeDeleteSegment(context.Background(), newSegment.Segment)
	if err != nil {
		log.Error("failed to execute query", logger.Err(err))
		render.JSON(w, r, Error("failed to execute query"))
		return
	}
	response := SegmentResponse{
		ResponseStatus: OK(),
		Segment:        newSegment.Segment,
	}
	log.Info("query successfully executed", slog.Any("request", response))
	render.JSON(w, r, response)
	return
}

func (s *ServerAPI) HandleGetSegmentUsersInfo(w http.ResponseWriter, r *http.Request) {
	return
}

package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/vlasashk/user-segmentation/internal/model/logger"
	"log/slog"
	"net/http"
	"os"
	"strconv"
)

// HandleAddUser godoc
// @Summary Add a new user
// @Description Add a new user to the system
// @ID addUser
// @Accept  json
// @Produce  json
// @Param user body UserRequest true "User object to be added"
// @Success 201 {object} UserResponse "Successfully added user"
// @Failure 400 {object} ResponseStatus "Invalid input data"
// @Failure 409 {object} ResponseStatus "Query execution failure"
// @Router /user/new [post]
func (s *ServerAPI) HandleAddUser(w http.ResponseWriter, r *http.Request) {
	newUser := &UserRequest{}
	log := s.Log.With(
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	if err := render.DecodeJSON(r.Body, &newUser); err != nil {
		log.Error("failed to decode request body", logger.Err(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, Error("failed to decode request body"))
		return
	}
	log.Info("request body decoded", slog.Any("request", *newUser))
	if err := validator.New().Struct(newUser); err != nil {
		log.Error("wrong body structure", logger.Err(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, Error("wrong body structure"))
		return
	}
	id, err := s.Store.AddUser(context.Background(), newUser.User, log)
	if err != nil {
		render.Status(r, http.StatusConflict)
		render.JSON(w, r, Error(err.Error()))
		return
	}
	response := UserResponse{
		ResponseStatus: OK(),
		User:           newUser.User,
	}
	response.User.Id = id
	log.Info("query successfully executed", slog.Any("request", response))
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
	return
}

// HandleAddSegment godoc
// @Summary Add a new segment
// @Description Add a new segment to the system
// @ID addSegment
// @Accept  json
// @Produce  json
// @Param segment body SegmentRequest true "Segment object to be added"
// @Success 201 {object} SegmentResponse "Successfully added segment"
// @Failure 400 {object} ResponseStatus "Invalid input data"
// @Failure 409 {object} ResponseStatus "Query execution failure"
// @Router /segment/new [post]
func (s *ServerAPI) HandleAddSegment(w http.ResponseWriter, r *http.Request) {
	newSegment := &SegmentRequest{}
	log := s.Log.With(
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	if err := render.DecodeJSON(r.Body, &newSegment); err != nil {
		log.Error("failed to decode request body", logger.Err(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, Error("failed to decode request body"))
		return
	}
	log.Info("request body decoded", slog.Any("request", *newSegment))
	if err := validator.New().Struct(newSegment); err != nil {
		log.Error("wrong body structure", logger.Err(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, Error("wrong body structure"))
		return
	}
	id, err := s.Store.AddSegment(context.Background(), newSegment.Segment, log)
	if err != nil {
		render.Status(r, http.StatusConflict)
		render.JSON(w, r, Error(err.Error()))
		return
	}
	response := SegmentResponse{
		ResponseStatus: OK(),
		Segment:        newSegment.Segment,
	}
	response.Segment.Id = id
	log.Info("query successfully executed", slog.Any("request", response))
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
	return
}

// HandleAddUserToSegment godoc
// @Summary Add user ti a segment
// @Description Link user to segments
// @ID addUserToSegment
// @Accept  json
// @Produce  json
// @Param segment body UserSegmentRequest true "Segment object to be added"
// @Success 201 {object} UserSegmentRequest "Successfully linked segment to a user"
// @Failure 400 {object} ResponseStatus "Invalid input data"
// @Failure 409 {object} ResponseStatus "Query execution failure"
// @Router /user/addSegment [post]
func (s *ServerAPI) HandleAddUserToSegment(w http.ResponseWriter, r *http.Request) {
	newUserSegment := &UserSegmentRequest{}
	log := s.Log.With(
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	if err := render.DecodeJSON(r.Body, &newUserSegment); err != nil {
		log.Error("failed to decode request body", logger.Err(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, Error("failed to decode request body"))
		return
	}
	log.Info("request body decoded", slog.Any("request", *newUserSegment))
	if err := validator.New().Struct(newUserSegment); err != nil {
		log.Error("wrong body structure", logger.Err(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, Error("wrong body structure"))
		return
	}
	if len(newUserSegment.UserSegments.SegmentSlug) == 0 {
		log.Error("empty segment array")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, Error("empty segment array"))
		return
	}
	err := s.Store.AddUserToSegments(context.Background(), newUserSegment.UserSegments, log)
	if err != nil {
		render.Status(r, http.StatusConflict)
		render.JSON(w, r, Error(err.Error()))
		return
	}
	response := UserSegmentResponse{
		ResponseStatus: OK(),
		UserSegments:   newUserSegment.UserSegments,
	}
	log.Info("query successfully executed", slog.Any("request", response))
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
	return
}

// HandleGetUserSegmentsInfo godoc
// @Summary Get user's segments information
// @Description Get information about the segments a user belongs to
// @ID getUserSegmentsInfo
// @Produce  json
// @Param userID path string true "user ID to get list of segments for"
// @Success 200 {object} GetSegmentsResponse "Successfully retrieved user segments"
// @Failure 400 {object} ResponseStatus "Invalid input data"
// @Failure 409 {object} ResponseStatus "Query execution failure"
// @Router /user/segments/{userID} [get]
func (s *ServerAPI) HandleGetUserSegmentsInfo(w http.ResponseWriter, r *http.Request) {
	user := &UserRequest{}
	log := s.Log.With(
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	if uid, err := strconv.ParseUint(chi.URLParam(r, "userID"), 10, 64); err != nil {
		log.Error("failed to parse user ID", logger.Err(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, Error("failed to parse user ID"))
		return
	} else {
		user.UID = uid
	}
	log.Info("user ID parsed successfully", slog.Any("request", *user))
	if err := validator.New().Struct(user); err != nil {
		log.Error("wrong body structure", logger.Err(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, Error("wrong body structure"))
		return
	}
	segments, err := s.Store.GetUserSegmentsInfo(context.Background(), user.User, log)
	if err != nil {
		render.Status(r, http.StatusConflict)
		render.JSON(w, r, Error(err.Error()))
		return
	}
	response := GetSegmentsResponse{
		ResponseStatus: OK(),
		UserID:         user.UID,
		SegmentSlug:    segments,
	}
	log.Info("query successfully executed", slog.Any("request", response))
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
	return
}

// HandleDeleteUserFromSegment godoc
// @Summary Remove a user from one or more segments
// @Description Remove a user from one or more segments
// @ID deleteUserFromSegment
// @Accept  json
// @Produce  json
// @Param userSegments body UserSegmentRequest true "User and segment association"
// @Success 200 {object} UserSegmentResponse "Successfully removed user from segment"
// @Failure 400 {object} ResponseStatus "Invalid input data"
// @Failure 409 {object} ResponseStatus "Query execution failure"
// @Router /user/segments [delete]
func (s *ServerAPI) HandleDeleteUserFromSegment(w http.ResponseWriter, r *http.Request) {
	newUserSegment := &UserSegmentRequest{}
	log := s.Log.With(
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	if err := render.DecodeJSON(r.Body, &newUserSegment); err != nil {
		log.Error("failed to decode request body", logger.Err(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, Error("failed to decode request body"))
		return
	}
	log.Info("request body decoded", slog.Any("request", *newUserSegment))
	if err := validator.New().Struct(newUserSegment); err != nil {
		log.Error("wrong body structure", logger.Err(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, Error("wrong body structure"))
		return
	}
	if len(newUserSegment.UserSegments.SegmentSlug) == 0 {
		log.Error("empty segment array")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, Error("empty segment array"))
		return
	}
	err := s.Store.DeleteUserFromSegments(context.Background(), newUserSegment.UserSegments, log)
	if err != nil {
		render.Status(r, http.StatusConflict)
		render.JSON(w, r, Error(err.Error()))
		return
	}
	response := UserSegmentResponse{
		ResponseStatus: OK(),
		UserSegments:   newUserSegment.UserSegments,
	}
	log.Info("query successfully executed", slog.Any("request", response))
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
	return
}

// HandleCascadeDeleteSegment godoc
// @Summary Cascade delete a segment
// @Description Cascade delete a segment and remove associated users
// @ID cascadeDeleteSegment
// @Accept  json
// @Produce  json
// @Param segment body SegmentRequest true "Segment object to delete"
// @Success 200 {object} SegmentResponse "Successfully deleted segment"
// @Failure 400 {object} ResponseStatus "Invalid input data"
// @Failure 409 {object} ResponseStatus "Query execution failure"
// @Router /segment/remove [delete]
func (s *ServerAPI) HandleCascadeDeleteSegment(w http.ResponseWriter, r *http.Request) {
	newSegment := &SegmentRequest{}
	log := s.Log.With(
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	if err := render.DecodeJSON(r.Body, &newSegment); err != nil {
		log.Error("failed to decode request body", logger.Err(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, Error("failed to decode request body"))
		return
	}
	log.Info("request body decoded", slog.Any("request", *newSegment))
	if err := validator.New().Struct(newSegment); err != nil {
		log.Error("wrong body structure", logger.Err(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, Error("wrong body structure"))
		return
	}
	err := s.Store.CascadeDeleteSegment(context.Background(), newSegment.Segment, log)
	if err != nil {
		render.Status(r, http.StatusConflict)
		render.JSON(w, r, Error(err.Error()))
		return
	}
	response := SegmentResponse{
		ResponseStatus: OK(),
		Segment:        newSegment.Segment,
	}
	log.Info("query successfully executed", slog.Any("request", response))
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
	return
}

// HandleGetSegmentUsersInfo godoc
// @Summary Get users of a segment
// @Description Get a list of users belonging to a specific segment
// @ID getSegmentUsersInfo
// @Produce  json
// @Param segmentName path string true "segment Name to get list of its users"
// @Success 200 {object} GetUsersResponse "Successfully retrieved segment users"
// @Failure 400 {object} ResponseStatus "Invalid input data"
// @Failure 409 {object} ResponseStatus "Query execution failure"
// @Router /segment/users/{segmentName} [get]
func (s *ServerAPI) HandleGetSegmentUsersInfo(w http.ResponseWriter, r *http.Request) {
	segment := &SegmentRequest{}
	log := s.Log.With(
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	segment.Slug = chi.URLParam(r, "segmentName")
	log.Info("segment name received", slog.Any("request", *segment))
	if err := validator.New().Struct(segment); err != nil {
		log.Error("wrong body structure", logger.Err(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, Error("wrong body structure"))
		return
	}
	users, err := s.Store.GetSegmentUsersInfo(context.Background(), segment.Segment, log)
	if err != nil {
		render.Status(r, http.StatusConflict)
		render.JSON(w, r, Error(err.Error()))
		return
	}
	response := GetUsersResponse{
		ResponseStatus: OK(),
		SegmentSlug:    segment.Slug,
		UserIDs:        users,
	}
	log.Info("query successfully executed", slog.Any("request", response))
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
	return
}

// HandleCsvReport godoc
// @Summary Generate CSV report
// @Description Generate a CSV report for a specific month and year
// @ID generateCsvReport
// @Accept  json
// @Produce  json
// @Param csvReport body CsvReportRequest true "CSV report request"
// @Success 200 {object} CsvReportResponse "Successfully generated CSV report"
// @Failure 400 {object} ResponseStatus "Invalid input data"
// @Failure 409 {object} ResponseStatus "Query execution failure"
// @Router /report [post]
func (s *ServerAPI) HandleCsvReport(w http.ResponseWriter, r *http.Request) {
	dates := &CsvReportRequest{}
	log := s.Log.With(
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	if err := render.DecodeJSON(r.Body, &dates); err != nil {
		log.Error("failed to decode request body", logger.Err(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, Error("failed to decode request body"))
		return
	}
	log.Info("request body decoded", slog.Any("request", *dates))
	if err := validator.New().Struct(dates); err != nil {
		log.Error("wrong body structure", logger.Err(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, Error("wrong body structure"))
		return
	}
	if !(dates.Month > 0 && dates.Month < 13) {
		log.Error("wrong body structure", logger.Err(fmt.Errorf("wrong month format")))
		render.Status(r, http.StatusConflict)
		render.JSON(w, r, Error("wrong month format"))
		return
	}
	if err := s.Store.CsvHistoryReport(context.Background(), dates.CsvReport, log); err != nil {
		render.Status(r, http.StatusConflict)
		render.JSON(w, r, Error(err.Error()))
		return
	}
	url := fmt.Sprintf("http://localhost:%s/report/%s_%d_%d.csv", os.Getenv("PORT"), "report", dates.Year, dates.Month)
	response := CsvReportResponse{
		ResponseStatus: OK(),
		CsvUrl:         url,
	}
	log.Info("query successfully executed", slog.Any("request", response))
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
	return
}

// HandleDownloadCsv godoc
// @Summary Download CSV report
// @Description Download a previously generated CSV report
// @ID downloadCsvReport
// @Produce  text/csv
// @Param fileName path string true "CSV file name to download"
// @Success 200 {file} file "CSV file for download"
// @Failure 400 {object} ResponseStatus "Invalid file name"
// @Router /report/{fileName} [get]
func (s *ServerAPI) HandleDownloadCsv(w http.ResponseWriter, r *http.Request) {
	log := s.Log.With(
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	fileName := fmt.Sprintf("%s%s.csv", os.Getenv("CSV_PATH"), chi.URLParam(r, "fileName"))
	if len(fileName) < 1 {
		log.Error("file name is empty", logger.Err(fmt.Errorf("file name is empty")))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, Error("file name is empty"))
		return
	}
	log.Info("file name acquired", slog.Any("request", fileName))
	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		log.Error("file doesn't exist", logger.Err(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, Error("file doesn't exist"))
		return
	}
	log.Info("file successfully found", slog.Any("request", fileName))
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	http.ServeFile(w, r, fileName)
	return
}

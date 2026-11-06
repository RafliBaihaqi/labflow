package handler

import (
	"errors"
	"labflow/dto"
	"labflow/lifecycle"
	"labflow/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReservationHandler struct {
	rsvc *service.ReservationService
}

func NewReservationHandler(rsvc *service.ReservationService) *ReservationHandler {
	return &ReservationHandler{rsvc: rsvc}
}

func (rh *ReservationHandler) CreateReservation(c *gin.Context) {
	var req dto.CreateReservationRequest
	if err := dto.BindStrictJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized", Message: "user id not found"})
		return
	}

	tests := make([]service.TestItemParams, len(req.Tests))
	for i, t := range req.Tests {
		tests[i] = service.TestItemParams{
			TestListID: t.TestListID,
			Method:     t.Method,
			Criteria:   t.Criteria,
			TestUnit:   t.TestUnit,
		}
	}

	resp, err := rh.rsvc.CreateReservation(c.Request.Context(), userID.(int), service.CreateReservationParams{
		SampleName:          req.SampleName,
		SampleAmount:        req.SampleAmount,
		ReceivedNote:        req.ReceivedNote,
		SampleResidue:       req.SampleResidue,
		SpecialInstructions: req.SpecialInstructions,
		SampleType:          req.SampleType,
		Note:                req.Note,
		BatchNumber:         req.BatchNumber,
		SampleTakenDate:     req.SampleTakenDate,
		SampleTakenBy:       req.SampleTakenBy,
		Tests:               tests,
	})
	if err != nil {
		status, body := mapReservationError(err)
		c.JSON(status, body)
		return
	}

	respTests := make([]dto.CreateTestRequest, len(resp.Tests))
	for i, t := range resp.Tests {
		respTests[i] = dto.CreateTestRequest{
			TestListID: t.TestListID,
			Method:     t.Method,
			Criteria:   t.Criteria,
			TestUnit:   t.TestUnit,
		}
	}

	c.JSON(http.StatusCreated, dto.CreateReservationResponse{
		ID:                  resp.ID,
		Xid:                 resp.Xid,
		SampleName:          resp.SampleName,
		SampleAmount:        resp.SampleAmount,
		SpecialInstructions: resp.SpecialInstructions,
		SampleType:          resp.SampleType,
		BatchNumber:         resp.BatchNumber,
		SampleTakenDate:     resp.SampleTakenDate,
		SampleTakenBy:       resp.SampleTakenBy,
		TestRequestNumber:   resp.TestRequestNumber,
		Message:             "Reservation created",
		Tests:               respTests,
	})
}

func (rh *ReservationHandler) UpdateReservationsByXid(c *gin.Context) {
	var req dto.UpdateReservationRequest
	if err := dto.BindStrictJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized", Message: "user id not found"})
		return
	}

	xid := c.Param("xid")
	xidParse, err := uuid.Parse(xid)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid xid", Message: err.Error()})
		return
	}
	if _, err := rh.rsvc.UpdateReservationByXid(c.Request.Context(), userID.(int), xidParse, service.UpdateReservationParams{
		POURL:       req.POUrl,
		PONumber:    req.PONumber,
		SamplePhoto: req.SamplePhoto,
	}); err != nil {
		status, body := mapReservationError(err)
		c.JSON(status, body)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Reservation updated"})
}

func (rh *ReservationHandler) ListReservations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	status := c.QueryArray("status")
	sortBy := c.Query("sort_by")
	sortOrder := c.Query("sort_order")

	resp, err := rh.rsvc.ListReservations(c.Request.Context(), service.ListParams{
		Search:    search,
		Status:    status,
		Page:      page,
		Limit:     limit,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	})
	if err != nil {
		status, body := mapReservationError(err)
		c.JSON(status, body)
		return
	}

	var data []dto.ReservationResponse
	for _, r := range resp.Data {
		data = append(data, toReservationResponse(r, "Reservations Fetched Successfully"))
	}

	c.JSON(http.StatusOK, dto.PaginatedReservationResponse{
		Data:       data,
		Total:      resp.Total,
		Page:       resp.Page,
		Limit:      resp.Limit,
		TotalPages: resp.TotalPages,
	})
}

func (rh *ReservationHandler) GetReservationStats(c *gin.Context) {
	resp, err := rh.rsvc.GetReservationStats(c.Request.Context())
	if err != nil {
		status, body := mapReservationError(err)
		c.JSON(status, body)
		return
	}

	c.JSON(http.StatusOK, dto.ReservationStatsResponse{
		Total:          resp.Total,
		Pending:        resp.Pending,
		SampleReceived: resp.SampleReceived,
		TestOnProgress: resp.TestOnProgress,
		TestCompleted:  resp.TestCompleted,
		Rejected:       resp.Rejected,
	})
}

func (rh *ReservationHandler) GetReservationByXid(c *gin.Context) {
	xid := c.Param("xid")
	xidParse, err := uuid.Parse(xid)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid xid", Message: err.Error()})
		return
	}

	resp, err := rh.rsvc.GetReservationByXid(c.Request.Context(), xidParse)
	if err != nil {
		status, body := mapReservationError(err)
		c.JSON(status, body)
		return
	}

	c.JSON(http.StatusOK, toReservationResponse(resp, "Reservation fetched"))
}

func (rh *ReservationHandler) UpdateReservationStatus(c *gin.Context) {
	var req dto.UpdateStatusReq
	if err := dto.BindStrictJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized", Message: "user id not found"})
		return
	}

	xid := c.Param("xid")
	xidParse, err := uuid.Parse(xid)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid xid", Message: err.Error()})
		return
	}
	if _, err := rh.rsvc.UpdateReservationStatus(c.Request.Context(), userID.(int), xidParse, service.UpdateStatusParams{
		Status:      req.Status,
		SamplePhoto: req.SamplePhoto,
	}); err != nil {
		status, body := mapReservationError(err)
		c.JSON(status, body)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Reservation status updated"})
}

func toReservationResponse(r service.ReservationResult, message string) dto.ReservationResponse {
	tests := make([]dto.TestResponse, 0, len(r.Tests))
	users := dto.UserResponse{
		ID:          r.User.ID,
		XID:         r.User.XID,
		Email:       r.User.Email,
		PhoneNumber: r.User.PhoneNumber,
		Name:        r.User.Name,
		Company:     r.User.Company,
		Address:     r.User.Address,
		Avatar:      r.User.Avatar,
	}
	for _, t := range r.Tests {
		tests = append(tests, dto.TestResponse{
			ID:            t.ID,
			XID:           t.XID,
			ReservationID: t.ReservationID,
			TestListId:    t.TestListId,
			Method:        t.Method,
			Criteria:      t.Criteria,
			Status:        t.Status,
			TestUnit:      t.TestUnit,
			TestResult:    t.TestResult,
			TestSummary:   t.TestSummary,
			ApprovedDate:  t.ApprovedDate,
			Notes:         t.Notes,
			CreatedAt:     t.CreatedAt,
			UpdatedAt:     t.UpdatedAt,
			TestList:      t.TestList,
		})
	}

	return dto.ReservationResponse{
		ID:                  r.ID,
		Xid:                 r.Xid,
		UserID:              r.UserID,
		ReservationDate:     r.ReservationDate,
		Status:              r.Status,
		SampleName:          r.SampleName,
		SampleAmount:        r.SampleAmount,
		ReceivedNote:        r.ReceivedNote,
		TestRequestNumber:   r.TestRequestNumber,
		SampleResidue:       r.SampleResidue,
		SpecialInstructions: r.SpecialInstructions,
		POURL:               r.POURL,
		PONumber:            r.PONumber,
		SampleType:          r.SampleType,
		FeedbackRequired:    r.FeedbackRequired,
		MRApproved:          r.MRApproved,
		HeadApproved:        r.HeadApproved,
		MRApprovalDate:      r.MRApprovalDate,
		HeadApprovalDate:    r.HeadApprovalDate,
		SamplePhoto:         r.SamplePhoto,
		Note:                r.Note,
		StartTestDate:       r.StartTestDate,
		EndTestDate:         r.EndTestDate,
		SAPNotes:            r.SAPNotes,
		BatchNumber:         r.BatchNumber,
		LeadTime:            r.LeadTime,
		SampleTakenDate:     r.SampleTakenDate,
		SampleTakenBy:       r.SampleTakenBy,
		IsKan:               r.IsKan,
		CreatedAt:           r.CreatedAt,
		UpdatedAt:           r.UpdatedAt,
		UpdatedBy:           r.UpdatedBy,
		Message:             message,
		Tests:               tests,
		User:                users,
	}
}

func mapReservationError(err error) (int, dto.ErrorResponse) {
	switch {
	case errors.Is(err, service.ErrFailedCreate):
		return errResp(http.StatusBadRequest, "failed_create", err.Error())
	case errors.Is(err, service.ErrFailedFetch):
		return errResp(http.StatusBadRequest, "failed_fetch", err.Error())
	case errors.Is(err, service.ErrFailedUpdate):
		return errResp(http.StatusBadRequest, "failed_update", err.Error())
	case errors.Is(err, lifecycle.ErrInvalidTransition):
		return errResp(http.StatusBadRequest, "invalid_status_update", err.Error())
	case errors.Is(err, lifecycle.ErrSamplePhotoRequired):
		return errResp(http.StatusBadRequest, "sample_photo_required", err.Error())
	}
	return errResp(http.StatusInternalServerError, "internal_error", "failed to process reservation")
}

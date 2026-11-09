package handler

import (
	"errors"
	"labflow/dto"
	"labflow/httpx"
	"labflow/service"
	"net/http"

	"github.com/google/uuid"
)

type TestHandler struct {
	tsvc *service.TestService
}

func NewTestHandler(tsvc *service.TestService) *TestHandler {
	return &TestHandler{tsvc: tsvc}
}

func (th *TestHandler) UpdateTestByXid(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateTestRequest
	if err := dto.BindStrictJSON(r, &req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	userID, exists := httpx.UserID(r.Context())
	if !exists {
		httpx.WriteJSON(w, http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized", Message: "user id not found"})
		return
	}

	xid := r.PathValue("xid")
	xidParse, err := uuid.Parse(xid)
	if err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid xid", Message: err.Error()})
		return
	}
	if _, err := th.tsvc.UpdateTestByXid(r.Context(), xidParse, userID, service.UpdateTestParams{
		TestResult:  req.TestResult,
		TestSummary: req.TestSummary,
		Notes:       req.Notes,
	}); err != nil {
		status, body := mapTestError(err)
		httpx.WriteJSON(w, status, body)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, dto.SuccessResponse{Message: "Test updated"})
}

func (th *TestHandler) UpdateTestApproval(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateTestRequest
	if err := dto.BindStrictJSON(r, &req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	userID, exists := httpx.UserID(r.Context())
	if !exists {
		httpx.WriteJSON(w, http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized", Message: "user id not found"})
		return
	}

	xid := r.PathValue("xid")
	xidParse, err := uuid.Parse(xid)
	if err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid xid", Message: err.Error()})
		return
	}
	if _, err := th.tsvc.UpdateIsApproved(r.Context(), xidParse, userID, service.UpdateTestParams{
		IsApproved: req.IsApproved,
	}); err != nil {
		status, body := mapTestError(err)
		httpx.WriteJSON(w, status, body)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, dto.SuccessResponse{Message: "Test updated"})
}

func mapTestError(err error) (int, dto.ErrorResponse) {
	switch {
	case errors.Is(err, service.ErrUnauthorizedRole):
		return errResp(http.StatusForbidden, "forbidden", err.Error())
	case errors.Is(err, service.ErrFailedFetch):
		return errResp(http.StatusNotFound, "not_found", "test not found")
	default:
		return errResp(http.StatusInternalServerError, "internal error", "failed to update test")
	}
}

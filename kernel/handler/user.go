package handler

import (
	"context"
	"errors"
	"labflow/dto"
	"labflow/httpx"
	"labflow/service"
	"net/http"

	"github.com/google/uuid"
)

type userService interface {
	UpdateUserByXid(ctx context.Context, updatedBy int, xid uuid.UUID, req service.UpdateUserParams) (service.UserResult, error)
}

type UserHandler struct {
	svc userService
}

func NewUserHandler(svc userService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, exists := httpx.UserID(r.Context())
	if !exists {
		httpx.WriteJSON(w, http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized", Message: "user not found"})
		return
	}

	xid, err := uuid.Parse(r.PathValue("xid"))
	if err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "bad_request", Message: "invalid xid"})
		return
	}

	var req dto.UpdateUserRequest
	if err := dto.BindStrictJSON(r, &req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	resp, err := h.svc.UpdateUserByXid(r.Context(), userID, xid, service.UpdateUserParams{
		PhoneNumber: req.PhoneNumber,
		Company:     req.Company,
		Address:     req.Address,
		Avatar:      req.Avatar,
	})
	if err != nil {
		status, body := mapUserError(err)
		httpx.WriteJSON(w, status, body)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, dto.UpdateUserResponse{
		Xid:         resp.Xid,
		PhoneNumber: resp.PhoneNumber,
		Company:     resp.Company,
		Address:     resp.Address,
		Avatar:      resp.Avatar,
		Message:     resp.Message,
	})
}

func mapUserError(err error) (int, dto.ErrorResponse) {
	if errors.Is(err, service.ErrFailedUpdate) {
		return errResp(http.StatusBadRequest, "bad_request", err.Error())
	}
	return errResp(http.StatusInternalServerError, "internal_error", "failed to update user")
}

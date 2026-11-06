package handler

import (
	"context"
	"errors"
	"labflow/dto"
	"labflow/service"
	"net/http"

	"github.com/gin-gonic/gin"
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

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized", Message: "user not found"})
		return
	}

	xid, err := uuid.Parse(c.Param("xid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "bad_request", Message: "invalid xid"})
		return
	}

	var req dto.UpdateUserRequest
	if err := dto.BindStrictJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	resp, err := h.svc.UpdateUserByXid(c.Request.Context(), userID.(int), xid, service.UpdateUserParams{
		PhoneNumber: req.PhoneNumber,
		Company:     req.Company,
		Address:     req.Address,
		Avatar:      req.Avatar,
	})
	if err != nil {
		status, body := mapUserError(err)
		c.JSON(status, body)
		return
	}

	c.JSON(http.StatusOK, dto.UpdateUserResponse{
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

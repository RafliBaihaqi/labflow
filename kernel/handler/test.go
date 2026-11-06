package handler

import (
	"errors"
	"labflow/dto"
	"labflow/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TestHandler struct {
	tsvc *service.TestService
}

func NewTestHandler(tsvc *service.TestService) *TestHandler {
	return &TestHandler{tsvc: tsvc}
}

func (th *TestHandler) UpdateTestByXid(c *gin.Context) {
	var req dto.UpdateTestRequest
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
	if _, err := th.tsvc.UpdateTestByXid(c.Request.Context(), xidParse, userID.(int), service.UpdateTestParams{
		TestResult:  req.TestResult,
		TestSummary: req.TestSummary,
		Notes:       req.Notes,
	}); err != nil {
		status, body := mapTestError(err)
		c.JSON(status, body)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Test updated"})
}

func (th *TestHandler) UpdateTestApproval(c *gin.Context) {
	var req dto.UpdateTestRequest
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
	if _, err := th.tsvc.UpdateIsApproved(c.Request.Context(), xidParse, userID.(int), service.UpdateTestParams{
		IsApproved: req.IsApproved,
	}); err != nil {
		status, body := mapTestError(err)
		c.JSON(status, body)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Test updated"})
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

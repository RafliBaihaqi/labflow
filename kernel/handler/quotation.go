package handler

import (
	"errors"
	"labflow/apperr"
	"labflow/dto"
	"labflow/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type QuotationHandler struct {
	qsvc *service.QuotationService
}

func NewQuotationHandler(qsvc *service.QuotationService) *QuotationHandler {
	return &QuotationHandler{qsvc: qsvc}
}

func (qh *QuotationHandler) CreateQuotation(c *gin.Context) {
	var req dto.CreateQuotationRequest
	if err := dto.BindStrictJSON(c, &req); err != nil {
		c.JSON(400, dto.ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	quo, err := qh.qsvc.CreateQuotation(c.Request.Context(), service.CreateQuotationParams{
		Email:       req.Email,
		CompanyName: req.CompanyName,
		PhoneNumber: req.PhoneNumber,
		Message:     req.Message,
	})
	if err != nil {
		c.JSON(500, dto.ErrorResponse{Error: "failed to create quotation", Message: err.Error()})
		return
	}

	c.JSON(201, dto.CreateQuotationResponse{
		XID:         quo.XID,
		Email:       quo.Email,
		CompanyName: quo.CompanyName,
		PhoneNumber: quo.PhoneNumber,
		Message:     quo.Message,
	})
}

func (qh *QuotationHandler) ListQuotations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	resp, err := qh.qsvc.ListQuotations(c.Request.Context(), search, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error", Message: err.Error()})
		return
	}

	data := make([]dto.QuotationResponse, 0, len(resp.Data))
	for _, q := range resp.Data {
		data = append(data, dto.QuotationResponse{
			XID:         q.XID,
			Email:       q.Email,
			CompanyName: q.CompanyName,
			PhoneNumber: q.PhoneNumber,
			Message:     q.Message,
			Status:      q.Status,
		})
	}

	c.JSON(http.StatusOK, dto.PaginatedQuotationResponse{
		Data:       data,
		Total:      resp.Total,
		Page:       resp.Page,
		Limit:      resp.Limit,
		TotalPages: resp.TotalPages,
	})
}

func (qh *QuotationHandler) UpdateQuotationStatus(c *gin.Context) {
	xid, err := uuid.Parse(c.Param("xid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid xid", Message: err.Error()})
		return
	}

	var req dto.UpdateQuotationStatusRequest
	if err := dto.BindStrictJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	quo, err := qh.qsvc.UpdateQuotationStatus(c.Request.Context(), xid, req.Status)
	if err != nil {
		status, body := mapQuotationError(err)
		c.JSON(status, body)
		return
	}

	c.JSON(http.StatusOK, dto.QuotationResponse{
		XID:         quo.XID,
		Email:       quo.Email,
		CompanyName: quo.CompanyName,
		PhoneNumber: quo.PhoneNumber,
		Message:     quo.Message,
		Status:      quo.Status,
	})
}

func mapQuotationError(err error) (int, dto.ErrorResponse) {
	switch {
	case errors.Is(err, service.ErrInvalidQuotationStatus):
		return errResp(http.StatusBadRequest, "invalid_status", err.Error())
	case errors.Is(err, apperr.ErrNotFound):
		return errResp(http.StatusNotFound, "not_found", err.Error())
	default:
		return errResp(http.StatusInternalServerError, "internal_error", "failed to process quotation")
	}
}

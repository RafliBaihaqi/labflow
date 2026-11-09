package handler

import (
	"errors"
	"labflow/apperr"
	"labflow/dto"
	"labflow/httpx"
	"labflow/service"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type QuotationHandler struct {
	qsvc *service.QuotationService
}

func NewQuotationHandler(qsvc *service.QuotationService) *QuotationHandler {
	return &QuotationHandler{qsvc: qsvc}
}

func (qh *QuotationHandler) CreateQuotation(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateQuotationRequest
	if err := dto.BindStrictJSON(r, &req); err != nil {
		httpx.WriteJSON(w, 400, dto.ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	quo, err := qh.qsvc.CreateQuotation(r.Context(), service.CreateQuotationParams{
		Email:       req.Email,
		CompanyName: req.CompanyName,
		PhoneNumber: req.PhoneNumber,
		Message:     req.Message,
	})
	if err != nil {
		httpx.WriteJSON(w, 500, dto.ErrorResponse{Error: "failed to create quotation", Message: err.Error()})
		return
	}

	httpx.WriteJSON(w, 201, dto.CreateQuotationResponse{
		XID:         quo.XID,
		Email:       quo.Email,
		CompanyName: quo.CompanyName,
		PhoneNumber: quo.PhoneNumber,
		Message:     quo.Message,
	})
}

func (qh *QuotationHandler) ListQuotations(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(httpx.QueryDefault(r, "page", "1"))
	limit, _ := strconv.Atoi(httpx.QueryDefault(r, "limit", "10"))
	search := r.URL.Query().Get("search")

	resp, err := qh.qsvc.ListQuotations(r.Context(), search, page, limit)
	if err != nil {
		httpx.WriteJSON(w, http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error", Message: err.Error()})
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

	httpx.WriteJSON(w, http.StatusOK, dto.PaginatedQuotationResponse{
		Data:       data,
		Total:      resp.Total,
		Page:       resp.Page,
		Limit:      resp.Limit,
		TotalPages: resp.TotalPages,
	})
}

func (qh *QuotationHandler) UpdateQuotationStatus(w http.ResponseWriter, r *http.Request) {
	xid, err := uuid.Parse(r.PathValue("xid"))
	if err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid xid", Message: err.Error()})
		return
	}

	var req dto.UpdateQuotationStatusRequest
	if err := dto.BindStrictJSON(r, &req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	quo, err := qh.qsvc.UpdateQuotationStatus(r.Context(), xid, req.Status)
	if err != nil {
		status, body := mapQuotationError(err)
		httpx.WriteJSON(w, status, body)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, dto.QuotationResponse{
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

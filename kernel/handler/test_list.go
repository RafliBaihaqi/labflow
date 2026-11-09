package handler

import (
	"errors"
	"labflow/apperr"
	"labflow/dto"
	"labflow/httpx"
	"labflow/models"
	"labflow/service"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type TestListHandler struct {
	tlsvc *service.TestListService
}

func NewTestListHandler(tlsvc *service.TestListService) *TestListHandler {
	return &TestListHandler{tlsvc: tlsvc}
}

func (tlh *TestListHandler) GetTestListByCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := r.PathValue("categoryid")
	catNum, err := strconv.Atoi(categoryID)
	if err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid category id", Message: err.Error()})
		return
	}
	results, err := tlh.tlsvc.GetTestListByCategory(r.Context(), catNum)
	if err != nil {
		status, body := mapTestListError(err)
		httpx.WriteJSON(w, status, body)
		return
	}

	resp := make([]dto.TestCatalogResponse, 0, len(results))
	for _, entry := range results {
		resp = append(resp, dto.TestCatalogResponse{
			ID:       entry.ID,
			Name:     entry.Name,
			Category: entry.Category,
		})
	}

	httpx.WriteJSON(w, http.StatusOK, resp)
}

func (tlh *TestListHandler) ListTestLists(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(httpx.QueryDefault(r, "page", "1"))
	limit, _ := strconv.Atoi(httpx.QueryDefault(r, "limit", "10"))
	search := r.URL.Query().Get("search")
	var categoryID *int
	if raw := r.URL.Query().Get("category_id"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			httpx.WriteJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid category id", Message: err.Error()})
			return
		}
		categoryID = &parsed
	}

	resp, err := tlh.tlsvc.ListTestLists(r.Context(), search, categoryID, page, limit)
	if err != nil {
		status, body := mapTestListAdminError(err)
		httpx.WriteJSON(w, status, body)
		return
	}

	data := make([]dto.TestListResponse, 0, len(resp.Data))
	for _, item := range resp.Data {
		data = append(data, toTestListResponse(&item))
	}

	httpx.WriteJSON(w, http.StatusOK, dto.PaginatedTestListResponse{
		Data:       data,
		Total:      resp.Total,
		Page:       resp.Page,
		Limit:      resp.Limit,
		TotalPages: resp.TotalPages,
	})
}

func (tlh *TestListHandler) CreateTestList(w http.ResponseWriter, r *http.Request) {
	var req dto.TestListInput
	if err := dto.BindStrictJSON(r, &req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	item, err := tlh.tlsvc.CreateTestList(r.Context(), serviceToModelTestList(req))
	if err != nil {
		status, body := mapTestListAdminError(err)
		httpx.WriteJSON(w, status, body)
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, toTestListResponse(item))
}

func (tlh *TestListHandler) UpdateTestList(w http.ResponseWriter, r *http.Request) {
	xid, err := uuid.Parse(r.PathValue("xid"))
	if err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid xid", Message: err.Error()})
		return
	}
	var req dto.TestListInput
	if err := dto.BindStrictJSON(r, &req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	item, err := tlh.tlsvc.UpdateTestList(r.Context(), xid, serviceToModelTestList(req))
	if err != nil {
		status, body := mapTestListAdminError(err)
		httpx.WriteJSON(w, status, body)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, toTestListResponse(item))
}

func (tlh *TestListHandler) DeactivateTestList(w http.ResponseWriter, r *http.Request) {
	xid, err := uuid.Parse(r.PathValue("xid"))
	if err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid xid", Message: err.Error()})
		return
	}
	item, err := tlh.tlsvc.DeactivateTestList(r.Context(), xid)
	if err != nil {
		status, body := mapTestListAdminError(err)
		httpx.WriteJSON(w, status, body)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, toTestListResponse(item))
}

func serviceToModelTestList(req dto.TestListInput) models.TestList {
	return models.TestList{
		CategoryID:         req.CategoryID,
		Value:              req.Value,
		Label:              req.Label,
		SampleAmountLiquid: req.SampleAmountLiquid,
		SampleAmountSolid:  req.SampleAmountSolid,
		SampleAmountPowder: req.SampleAmountPowder,
		Leadtime:           req.Leadtime,
	}
}

func toTestListResponse(item *models.TestList) dto.TestListResponse {
	if item == nil {
		return dto.TestListResponse{}
	}
	return dto.TestListResponse{
		XID:                item.XID.String(),
		CategoryID:         item.CategoryID,
		Value:              item.Value,
		Label:              item.Label,
		SampleAmountLiquid: item.SampleAmountLiquid,
		SampleAmountSolid:  item.SampleAmountSolid,
		SampleAmountPowder: item.SampleAmountPowder,
		Leadtime:           item.Leadtime,
		Version:            item.Version,
		IsCurrent:          item.IsCurrent,
		CategoryName:       item.Category.Name,
	}
}

func mapTestListAdminError(err error) (int, dto.ErrorResponse) {
	switch {
	case errors.Is(err, service.ErrInvalidTestList):
		return errResp(http.StatusBadRequest, "invalid_test_list", err.Error())
	case errors.Is(err, apperr.ErrNotFound):
		return errResp(http.StatusNotFound, "not_found", err.Error())
	default:
		return errResp(http.StatusInternalServerError, "internal_error", "failed to process test list")
	}
}

func mapTestListError(err error) (int, dto.ErrorResponse) {
	switch {
	case errors.Is(err, apperr.ErrNotFound):
		return errResp(http.StatusNotFound, "not_found", err.Error())
	default:
		return errResp(http.StatusInternalServerError, "internal_error", "failed to fetch test list")
	}
}

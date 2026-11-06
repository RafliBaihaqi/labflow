package handler

import (
	"errors"
	"labflow/apperr"
	"labflow/dto"
	"labflow/models"
	"labflow/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TestListHandler struct {
	tlsvc *service.TestListService
}

func NewTestListHandler(tlsvc *service.TestListService) *TestListHandler {
	return &TestListHandler{tlsvc: tlsvc}
}

func (tlh *TestListHandler) GetTestListByCategory(c *gin.Context) {
	categoryID := c.Param("categoryid")
	catNum, err := strconv.Atoi(categoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid category id", Message: err.Error()})
		return
	}
	results, err := tlh.tlsvc.GetTestListByCategory(c.Request.Context(), catNum)
	if err != nil {
		status, body := mapTestListError(err)
		c.JSON(status, body)
		return
	}

	resp := make([]dto.TestCatalogResponse, 0, len(results))
	for _, r := range results {
		resp = append(resp, dto.TestCatalogResponse{
			ID:       r.ID,
			Name:     r.Name,
			Category: r.Category,
		})
	}

	c.JSON(http.StatusOK, resp)
}

func (tlh *TestListHandler) ListTestLists(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	var categoryID *int
	if raw := c.Query("category_id"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid category id", Message: err.Error()})
			return
		}
		categoryID = &parsed
	}

	resp, err := tlh.tlsvc.ListTestLists(c.Request.Context(), search, categoryID, page, limit)
	if err != nil {
		status, body := mapTestListAdminError(err)
		c.JSON(status, body)
		return
	}

	data := make([]dto.TestListResponse, 0, len(resp.Data))
	for _, item := range resp.Data {
		data = append(data, toTestListResponse(&item))
	}

	c.JSON(http.StatusOK, dto.PaginatedTestListResponse{
		Data:       data,
		Total:      resp.Total,
		Page:       resp.Page,
		Limit:      resp.Limit,
		TotalPages: resp.TotalPages,
	})
}

func (tlh *TestListHandler) CreateTestList(c *gin.Context) {
	var req dto.TestListInput
	if err := dto.BindStrictJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	item, err := tlh.tlsvc.CreateTestList(c.Request.Context(), serviceToModelTestList(req))
	if err != nil {
		status, body := mapTestListAdminError(err)
		c.JSON(status, body)
		return
	}

	c.JSON(http.StatusCreated, toTestListResponse(item))
}

func (tlh *TestListHandler) UpdateTestList(c *gin.Context) {
	xid, err := uuid.Parse(c.Param("xid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid xid", Message: err.Error()})
		return
	}
	var req dto.TestListInput
	if err := dto.BindStrictJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	item, err := tlh.tlsvc.UpdateTestList(c.Request.Context(), xid, serviceToModelTestList(req))
	if err != nil {
		status, body := mapTestListAdminError(err)
		c.JSON(status, body)
		return
	}

	c.JSON(http.StatusOK, toTestListResponse(item))
}

func (tlh *TestListHandler) DeactivateTestList(c *gin.Context) {
	xid, err := uuid.Parse(c.Param("xid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid xid", Message: err.Error()})
		return
	}
	item, err := tlh.tlsvc.DeactivateTestList(c.Request.Context(), xid)
	if err != nil {
		status, body := mapTestListAdminError(err)
		c.JSON(status, body)
		return
	}
	c.JSON(http.StatusOK, toTestListResponse(item))
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

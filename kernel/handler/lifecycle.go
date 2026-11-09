package handler

import (
	"labflow/httpx"
	"labflow/lifecycle"
	"net/http"
)

type LifecycleHandler struct{}

func NewLifecycleHandler() *LifecycleHandler {
	return &LifecycleHandler{}
}

func (h *LifecycleHandler) GetTransitions(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"transitions":  lifecycle.Transitions(),
		"all_statuses": lifecycle.AllStatuses(),
	})
}

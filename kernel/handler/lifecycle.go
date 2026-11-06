package handler

import (
	"labflow/lifecycle"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LifecycleHandler struct{}

func NewLifecycleHandler() *LifecycleHandler {
	return &LifecycleHandler{}
}

func (h *LifecycleHandler) GetTransitions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"transitions":  lifecycle.Transitions(),
		"all_statuses": lifecycle.AllStatuses(),
	})
}

package handler

import (
	"labflow/dto"

	"github.com/gin-gonic/gin"
)

func errResp(status int, code, message string) (int, dto.ErrorResponse) {
	return status, dto.ErrorResponse{
		Error:   code,
		Message: message,
	}
}

func writeErr(c *gin.Context, status int, code, message string) {
	c.JSON(status, dto.ErrorResponse{
		Error:   code,
		Message: message,
	})
}

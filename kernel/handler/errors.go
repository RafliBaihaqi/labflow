package handler

import (
	"labflow/dto"
	"labflow/httpx"
	"net/http"
)

func errResp(status int, code, message string) (int, dto.ErrorResponse) {
	return status, dto.ErrorResponse{
		Error:   code,
		Message: message,
	}
}

func writeErr(w http.ResponseWriter, status int, code, message string) {
	httpx.WriteJSON(w, status, dto.ErrorResponse{
		Error:   code,
		Message: message,
	})
}

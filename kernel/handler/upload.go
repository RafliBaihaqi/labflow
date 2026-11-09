package handler

import (
	"context"
	"errors"
	"io"
	"labflow/dto"
	"labflow/httpx"
	"labflow/service"
	"net/http"
	"path/filepath"
	"strings"
)

type uploadService interface {
	Upload(ctx context.Context, dest service.UploadDestination, file io.Reader, size int64, contentType, extension string) (string, error)
}

type UploadHandler struct {
	svc uploadService
}

func NewUploadHandler(svc uploadService) *UploadHandler {
	return &UploadHandler{svc: svc}
}

func (h *UploadHandler) UploadSamplePhoto(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "bad_request", Message: "file is required"})
		return
	}
	defer file.Close()

	ext := strings.TrimPrefix(filepath.Ext(header.Filename), ".")
	if ext == "" {
		ext = "jpg"
	}

	objectName, err := h.svc.Upload(r.Context(), service.DestinationReservation, file, header.Size, header.Header.Get("Content-Type"), ext)
	if err != nil {
		status, body := mapUploadError(err)
		httpx.WriteJSON(w, status, body)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, dto.UploadResponse{
		Message:    "File uploaded successfully",
		ObjectName: objectName,
	})
}

func (h *UploadHandler) UploadCertificate(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "bad_request", Message: "file is required"})
		return
	}
	defer file.Close()

	ext := strings.TrimPrefix(filepath.Ext(header.Filename), ".")
	if ext == "" {
		ext = "pdf"
	}

	objectName, err := h.svc.Upload(r.Context(), service.DestinationCertificate, file, header.Size, header.Header.Get("Content-Type"), ext)
	if err != nil {
		status, body := mapUploadError(err)
		httpx.WriteJSON(w, status, body)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, dto.UploadResponse{
		Message:    "File uploaded successfully",
		ObjectName: objectName,
	})
}

func (h *UploadHandler) UploadUserAvatar(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "bad_request", Message: "file is required"})
		return
	}
	defer file.Close()

	ext := strings.TrimPrefix(filepath.Ext(header.Filename), ".")
	if ext == "" {
		ext = "jpg"
	}

	objectName, err := h.svc.Upload(r.Context(), service.DestinationUser, file, header.Size, header.Header.Get("Content-Type"), ext)
	if err != nil {
		status, body := mapUploadError(err)
		httpx.WriteJSON(w, status, body)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, dto.UploadResponse{
		Message:    "File uploaded successfully",
		ObjectName: objectName,
	})
}

func mapUploadError(err error) (int, dto.ErrorResponse) {
	if errors.Is(err, service.ErrUnknownDestination) {
		return errResp(http.StatusBadRequest, "bad_request", err.Error())
	}
	if errors.Is(err, service.ErrUploadFailed) {
		return errResp(http.StatusBadRequest, "bad_request", err.Error())
	}
	return errResp(http.StatusInternalServerError, "internal_error", "failed to upload file")
}

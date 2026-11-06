package handler

import (
	"context"
	"errors"
	"io"
	"labflow/dto"
	"labflow/service"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
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

func (h *UploadHandler) UploadSamplePhoto(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "bad_request", Message: "file is required"})
		return
	}
	defer file.Close()

	ext := strings.TrimPrefix(filepath.Ext(header.Filename), ".")
	if ext == "" {
		ext = "jpg"
	}

	objectName, err := h.svc.Upload(c.Request.Context(), service.DestinationReservation, file, header.Size, header.Header.Get("Content-Type"), ext)
	if err != nil {
		status, body := mapUploadError(err)
		c.JSON(status, body)
		return
	}

	c.JSON(http.StatusOK, dto.UploadResponse{
		Message:    "File uploaded successfully",
		ObjectName: objectName,
	})
}

func (h *UploadHandler) UploadCertificate(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "bad_request", Message: "file is required"})
		return
	}
	defer file.Close()

	ext := strings.TrimPrefix(filepath.Ext(header.Filename), ".")
	if ext == "" {
		ext = "pdf"
	}

	objectName, err := h.svc.Upload(c.Request.Context(), service.DestinationCertificate, file, header.Size, header.Header.Get("Content-Type"), ext)
	if err != nil {
		status, body := mapUploadError(err)
		c.JSON(status, body)
		return
	}

	c.JSON(http.StatusOK, dto.UploadResponse{
		Message:    "File uploaded successfully",
		ObjectName: objectName,
	})
}

func (h *UploadHandler) UploadUserAvatar(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "bad_request", Message: "file is required"})
		return
	}
	defer file.Close()

	ext := strings.TrimPrefix(filepath.Ext(header.Filename), ".")
	if ext == "" {
		ext = "jpg"
	}

	objectName, err := h.svc.Upload(c.Request.Context(), service.DestinationUser, file, header.Size, header.Header.Get("Content-Type"), ext)
	if err != nil {
		status, body := mapUploadError(err)
		c.JSON(status, body)
		return
	}

	c.JSON(http.StatusOK, dto.UploadResponse{
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

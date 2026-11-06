package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"labflow/dto"
	"labflow/service"

	"github.com/gin-gonic/gin"
)

// fileFetcher is the narrow dependency the handler needs to stream a stored file.
type fileFetcher interface {
	GetObject(ctx context.Context, bucketName, objectKey string) (*service.ObjectData, error)
}

type FileHandler struct {
	files             fileFetcher
	bucketReservation string
	bucketUser        string
	bucketCertificate string
}

func NewFileHandler(files fileFetcher, bucketReservation, bucketUser, bucketCertificate string) *FileHandler {
	return &FileHandler{
		files:             files,
		bucketReservation: bucketReservation,
		bucketUser:        bucketUser,
		bucketCertificate: bucketCertificate,
	}
}

func (h *FileHandler) ServeReservationFile(c *gin.Context) {
	h.serve(c, h.bucketReservation)
}

func (h *FileHandler) ServeUserFile(c *gin.Context) {
	h.serve(c, h.bucketUser)
}

func (h *FileHandler) ServeCertificateFile(c *gin.Context) {
	h.serve(c, h.bucketCertificate)
}

func (h *FileHandler) serve(c *gin.Context, bucketName string) {
	object := c.Param("object")
	if !validObjectName(object) {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "bad_request", Message: "invalid object name"})
		return
	}

	data, err := h.files.GetObject(c.Request.Context(), bucketName, object)
	if err != nil {
		if errors.Is(err, service.ErrObjectNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not_found", Message: "file not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error", Message: "failed to fetch file"})
		return
	}
	defer data.Body.Close()

	c.DataFromReader(http.StatusOK, data.Size, data.ContentType, data.Body, nil)
}

// validObjectName rejects empty names and anything that could traverse paths.
// Gin's :object param never contains a slash, but we defend against "..", path
// separators, and empty input explicitly.
func validObjectName(name string) bool {
	if name == "" {
		return false
	}
	if strings.Contains(name, "/") || strings.Contains(name, "\\") || strings.Contains(name, "..") {
		return false
	}
	return true
}

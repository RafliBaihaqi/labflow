package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"labflow/dto"
	"labflow/httpx"
	"labflow/service"
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

func (h *FileHandler) ServeReservationFile(w http.ResponseWriter, r *http.Request) {
	h.serve(w, r, h.bucketReservation)
}

func (h *FileHandler) ServeUserFile(w http.ResponseWriter, r *http.Request) {
	h.serve(w, r, h.bucketUser)
}

func (h *FileHandler) ServeCertificateFile(w http.ResponseWriter, r *http.Request) {
	h.serve(w, r, h.bucketCertificate)
}

func (h *FileHandler) serve(w http.ResponseWriter, r *http.Request, bucketName string) {
	object := r.PathValue("object")
	if !validObjectName(object) {
		httpx.WriteJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "bad_request", Message: "invalid object name"})
		return
	}

	data, err := h.files.GetObject(r.Context(), bucketName, object)
	if err != nil {
		if errors.Is(err, service.ErrObjectNotFound) {
			httpx.WriteJSON(w, http.StatusNotFound, dto.ErrorResponse{Error: "not_found", Message: "file not found"})
			return
		}
		httpx.WriteJSON(w, http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error", Message: "failed to fetch file"})
		return
	}
	defer data.Body.Close()

	w.Header().Set("Content-Type", data.ContentType)
	if data.Size >= 0 {
		w.Header().Set("Content-Length", strconv.FormatInt(data.Size, 10))
	}
	w.WriteHeader(http.StatusOK)
	_, _ = io.Copy(w, data.Body)
}

// validObjectName rejects empty names and anything that could traverse paths.
// A {object} path wildcard never contains a slash, but we defend against "..",
// path separators, and empty input explicitly.
func validObjectName(name string) bool {
	if name == "" {
		return false
	}
	if strings.Contains(name, "/") || strings.Contains(name, "\\") || strings.Contains(name, "..") {
		return false
	}
	return true
}

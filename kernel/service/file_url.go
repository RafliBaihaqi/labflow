package service

// File category path segments. These MUST match the route segments registered
// in server/serve.go under the /api/files group.
const (
	FileCategoryReservation = "reservations"
	FileCategoryUser        = "users"
	FileCategoryCertificate = "certificates"
)

// FileURL builds a relative, same-origin URL the browser can fetch through the
// /api/files endpoint. A nil or empty object name is passed through unchanged so
// callers can map optional photo/avatar fields directly.
func FileURL(category string, objectName *string) *string {
	if objectName == nil || *objectName == "" {
		return objectName
	}
	u := "/api/files/" + category + "/" + *objectName
	return &u
}

package dto

type PaginationParams struct {
	Page   int
	Limit  int
	Search string
	Status []string
}

type PaginatedReservationResponse struct {
	Data       []ReservationResponse `json:"data"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	Limit      int                   `json:"limit"`
	TotalPages int                   `json:"total_pages"`
}

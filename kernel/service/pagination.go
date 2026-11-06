package service

func normalPagination(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return page, limit
}

func calculateTotalPage(total int64, limit int) int {
	return int((total + int64(limit) - 1) / int64(limit))
}

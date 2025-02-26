package domain

type PaginationRequest struct {
	Page  int `form:"page,default=1" binding:"gte=1"`
	Limit int `form:"limit,default=10" binding:"gte=1,lte=100"`
}

type PaginatedResponse struct {
	Data        interface{} `json:"data"`
	TotalItems  int64       `json:"total_items"`
	CurrentPage int         `json:"current_page"`
	PerPage     int         `json:"per_page"`
	TotalPages  int         `json:"total_pages"`
}

func NewPaginatedResponse(data interface{}, totalItems int64, currentPage, perPage int) *PaginatedResponse {
	totalPages := (int(totalItems) + perPage - 1) / perPage
	return &PaginatedResponse{
		Data:        data,
		TotalItems:  int64(totalPages),
		CurrentPage: currentPage,
		PerPage:     perPage,
		TotalPages:  totalPages,
	}
}

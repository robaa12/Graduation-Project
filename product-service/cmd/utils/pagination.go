package utils

type PaginationQuery struct {
	Page     int
	PageSize int
	Sort     string
	Order    string
}

func NewPaginationQuery(page, pageSize int, sort, order string) *PaginationQuery {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if sort == "" {
		sort = "created_at"
	}
	if order == "" {
		order = "desc"
	}

	return &PaginationQuery{
		Page:     page,
		PageSize: pageSize,
		Sort:     sort,
		Order:    order,
	}
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

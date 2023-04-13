package pagination

type Pagination struct {
	Total       int `json:"total"`
	LastPage    int `json:"last_page"`
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	From        int `json:"from"`
	To          int `json:"to"`
}

type PaginationResult[T any] struct {
	Data       []T        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

var DEFAULT_PAGINATION = Pagination{
	Total:       0,
	PerPage:     10,
	From:        0,
	To:          10,
	CurrentPage: 1,
	LastPage:    1,
}

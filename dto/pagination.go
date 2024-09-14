package dto

type Next struct {
	Page int `json:"page"`
}

type Prev struct {
	Page int `json:"page"`
}

type PaginationResp[T any] struct {
	Total      int  `json:"total"`
	IsLoadMore bool `json:"isLoadMore"`
	Data       []T  `json:"data"`
	Next       Next `json:"next"`
	Prev       Prev `json:"prev"`
}

func ToPaginationResp[T any](data []T, page int, limit int, total int) PaginationResp[T] {
	var nextPage Next
	var prevPage Prev
	isLoadMore := false

	// FROM CLIENT
	startPage := (page - 1) * limit
	endPage := page * limit

	if endPage < total {
		nextPage.Page = page + 1
		isLoadMore = true
	} else {
		nextPage.Page = -1
		isLoadMore = false
	}

	if startPage > 0 {
		prevPage.Page = page - 1
	} else {
		prevPage.Page = -1
	}

	return PaginationResp[T]{
		Next:       nextPage,
		Prev:       prevPage,
		Total:      total,
		IsLoadMore: isLoadMore,
		Data:       data,
	}
}

package form

import (
	"strconv"

	"github.com/ranta0/rest-and-go/pagination"
)

type Pagination struct {
	Page    string `json:"page" form:"page" query:"page"`
	PerPage string `json:"per_page" form:"per_page" query:"per_page"`
}

func NewPaginatorFromRequest(form *Pagination, count int) *pagination.Paginator {
	page := setOrDefault(form.Page, 1)
	perPage := setOrDefault(form.PerPage, pagination.DefaultPageSize)

	return pagination.NewPaginator(page, perPage, count)
}

func setOrDefault(value string, defaultValue int) int {
	if result, err := strconv.Atoi(value); err == nil {
		return result
	}

	return defaultValue
}

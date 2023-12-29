package pagination

import (
	"strconv"

	"github.com/ranta0/rest-and-go/form"
)

var (
	defaultPageSize = 10
	maxPageSize     = 100
)

type Paginator struct {
	Page       int `json:"page,omitempty"`
	PerPage    int `json:"per_page,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
	TotalItems int `json:"total_items,omitempty"`
}

func NewPaginator(page, perPage, totalItems int) *Paginator {
	if perPage <= 0 {
		perPage = defaultPageSize
	}
	if perPage > maxPageSize {
		perPage = maxPageSize
	}
	pageCount := -1
	if totalItems >= 0 {
		pageCount = (totalItems + perPage - 1) / perPage
		if page > pageCount {
			page = pageCount
		}
	}
	if page < 1 {
		page = 1
	}

	return &Paginator{
		Page:       page,
		PerPage:    perPage,
		TotalPages: pageCount,
		TotalItems: totalItems,
	}
}

func (p *Paginator) Offset() int {
	return p.PerPage * (p.Page - 1)
}

func (p *Paginator) Limit() int {
	return p.PerPage
}

func NewFromRequest(form *form.Pagination, count int) *Paginator {
	page := parseInt(form.Page, 1)
	perPage := parseInt(form.PerPage, defaultPageSize)
	return NewPaginator(page, perPage, count)
}

func parseInt(value string, defaultValue int) int {
	if result, err := strconv.Atoi(value); err == nil {
		return result
	}
	return defaultValue
}

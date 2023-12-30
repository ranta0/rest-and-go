package pagination

var (
	DefaultPageSize = 10
	MaxPageSize     = 100
)

type Paginator struct {
	Page       int `json:"page,omitempty"`
	PerPage    int `json:"per_page,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
	TotalItems int `json:"total_items,omitempty"`
}

func NewPaginator(page, perPage, totalItems int) *Paginator {
	if perPage <= 0 {
		perPage = DefaultPageSize
	}
	if perPage > MaxPageSize {
		perPage = MaxPageSize
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

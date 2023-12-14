package form

type Pagination struct {
	Page    string `json:"page" form:"page" query:"page"`
	PerPage string `json:"per_page" form:"per_page" query:"per_page"`
}

package response

import (
	"fmt"

	"github.com/ranta0/rest-and-go/pagination"
)

type JSONStub struct {
	*pagination.Paginator
	Status  string            `json:"status,omitempty"`
	Message string            `json:"message,omitempty"`
	Data    interface{}       `json:"data,omitempty"`
	Links   map[string]string `json:"links,omitempty"`
}

func (j *JSONStub) addLink(rel string, href string) {
	j.Links[rel] = href
}

func (j *JSONStub) AddPaginationLinks(href string) {
	if j.Paginator == nil {
		return
	}

	var prevHref string
	if j.Page-1 > 0 {
		prevHref = href + fmt.Sprintf("?page=%d", (j.Page-1))
	}
	j.addLink("prev", prevHref)
	j.addLink("self", href+fmt.Sprintf("?page=%d", (j.Page)))
	j.addLink("next", href+fmt.Sprintf("?page=%d", (j.Page+1)))
	j.addLink("first", href+"?page=1")
	j.addLink("last", href+fmt.Sprintf("?page=%d", j.TotalPages))
}

func (j *JSONStub) AddSelfLink(href string) {
	j.addLink("self", href)
}

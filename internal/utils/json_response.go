package utils

import (
	"net/http"

	"github.com/go-chi/render"
)

type JSONResponse struct {
	Status  string            `json:"status,omitempty"`
	Message string            `json:"message,omitempty"`
	Data    interface{}       `json:"data,omitempty"`
	Links   map[string]string `json:"links,omitempty"`
}

func NewJSONResponse(v interface{}, status string, message string) *JSONResponse {
	return &JSONResponse{
		Data:    v,
		Status:  status,
		Message: message,
		Links: make(map[string]string),
	}
}

func SuccessJsonResponse(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	render.Status(r, code)
	render.JSON(w, r, NewJSONResponse(data, "success", ""))
}

func ErrorJsonResponse(w http.ResponseWriter, r *http.Request, code int, message string) {
	render.Status(r, code)
	render.JSON(w, r, NewJSONResponse(nil, "err", message))
}

func AddLink(response *JSONResponse, rel string, href string) {
	response.Links[rel] = href
}

package response

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/go-chi/render"
	"github.com/ranta0/rest-and-go/pagination"
)

func OK(w http.ResponseWriter, r *http.Request, code int, stub *JSONStub) {
	if stub.Data == nil {
		render.Status(r, code)
		render.JSON(w, r, stub)
		return
	}

	url := getFullURL(r)

	if isArray(stub.Data) {
		stub.AddPaginationLinks(url)
	} else {
		// add id resource to the url
		var param string
		value, exist := getInterfaceKeyAndValue(stub.Data, "PublicID")
		if exist {
			param = value.(string)
		}

		selfUrl := url
		if !strings.Contains(url, param) {
			selfUrl += param
		}
		stub.AddSelfLink(selfUrl)
	}

	render.Status(r, code)
	render.JSON(w, r, stub)
}

func Error(w http.ResponseWriter, r *http.Request, code int, message string) {
	render.Status(r, code)
	render.JSON(w, r, NewError(message))
}

func NewOK(v interface{}, pages *pagination.Paginator) *JSONStub {
	return &JSONStub{
		Data:      v,
		Status:    "success",
		Links:     make(map[string]string),
		Paginator: pages,
	}
}

func NewError(message string) *JSONStub {
	return &JSONStub{
		Message: message,
		Status:  "error",
	}
}

func getFullURL(r *http.Request) string {
	r.URL.RawQuery = ""

	scheme := r.URL.Scheme
	if scheme == "" {
		scheme = "http"
	}

	fullURL := scheme + "://" + r.Host + r.URL.String()
	return fullURL
}

func isArray(data interface{}) bool {
	value := reflect.ValueOf(data)

	if value.Kind() == reflect.Array || value.Kind() == reflect.Slice {
		return true
	}

	return false
}

func getInterfaceKeyAndValue(data interface{}, key string) (interface{}, bool) {
	structValue := reflect.ValueOf(data)
	if structValue.Kind() == reflect.Ptr {
		structValue = structValue.Elem()
	}

	typ := structValue.Type()
	if typ.Kind() != reflect.Struct {
		return nil, false
	}

	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)
		fieldName := typ.Field(i).Name

		if fieldName == key {
			return field.Interface(), true
		}
	}

	return nil, false
}

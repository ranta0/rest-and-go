package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	httpError "github.com/ranta0/rest-and-go/internal/error"
)

// Handles all kinds of content types
func Handler(r *http.Request, v interface{}) error {
	contentType := r.Header.Get("Content-Type")

	switch contentType {
	case "text/plain":
		return handleJSON(r, v)
	case "application/json":
		return handleJSON(r, v)
	case "application/x-www-form-urlencoded":
		return handleForm(r, v)
	default:
		return httpError.ErrContentType
	}
}

// Handles JSON content types
func handleJSON(r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return httpError.ErrPayload
	}

	return nil
}

// Handles FormEncoded content types
func handleForm(r *http.Request, v interface{}) error {
	err := parseFormEncode(r, v)
	if err != nil {
		return httpError.ErrPayload
	}

	return nil
}

// parseFormEncode grabs the info from the form notation in a type
func parseFormEncode(r *http.Request, v interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	formValues := r.Form

	structValue := reflect.ValueOf(v).Elem()
	structType := structValue.Type()

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := structValue.Field(i)

		// Ignore fields without the form tag
		formTag := field.Tag.Get("form")
		if formTag == "" {
			continue
		}

		// Ignore empty values
		value := formValues.Get(formTag)
		if value == "" {
			continue
		}

		switch fieldValue.Kind() {
		case reflect.String:
			fieldValue.SetString(value)
		case reflect.Int:
			intValue, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("could not convert %s to int: %v", formTag, err)
			}
			fieldValue.SetInt(int64(intValue))
		default:
			return fmt.Errorf("unsupported field type: %s", fieldValue.Kind())
		}
	}

	return nil
}

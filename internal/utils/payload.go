package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"

	httpErrors "github.com/ranta0/rest-and-go/internal/errors"
)

// Handles all kinds of content types
func HandlePayload(r *http.Request, v interface{}) error {
	contentType := r.Header.Get("Content-Type")

	switch contentType {
	case "text/plain":
		return handleJSON(r, v)
	case "application/json":
		return handleJSON(r, v)
	case "application/x-www-form-urlencoded":
		return handleForm(r, v)
	default:
		return httpErrors.ErrContentType
	}
}

// Handles all kinds of payload content types
func ValidatePayload(v interface{}) error {
	var validate = validator.New()

	err := validate.Struct(v)
	failedValidations := make([]string, 0)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			fieldName := e.StructField()
			jsonFieldName, found := getJSONFieldName(v, e.StructField())
			if found {
				fieldName = jsonFieldName
			}

			failedValidations = append(
				failedValidations,
				fmt.Sprintf("field '%s' failed validation for rule '%s'", fieldName, e.Tag()),
			)
		}
	}

	if len(failedValidations) == 0 {
		return nil
	}

	return fmt.Errorf("%s", strings.Join(failedValidations, ", "))
}

// Handles JSON content types
func handleJSON(r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return httpErrors.ErrPayload
	}

	return nil
}

// Handles FormEncoded content types
func handleForm(r *http.Request, v interface{}) error {
	err := parseFormEncode(r, v)
	if err != nil {
		return httpErrors.ErrPayload
	}

	return nil
}

// parseFormEncode grabs the info from the form notaion in a type
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

// getJSONFieldName retrieves the JSON field name for a given struct field
func getJSONFieldName(structInstance interface{}, fieldName string) (string, bool) {
	structValue := reflect.ValueOf(structInstance)
	if structValue.Kind() == reflect.Ptr {
		structValue = structValue.Elem()
	}

	typ := structValue.Type()
	if typ.Kind() != reflect.Struct {
		return "", false
	}

	field, found := typ.FieldByName(fieldName)
	if !found {
		return "", false
	}

	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return "", false
	}

	// Extract the JSON field name from the tag
	jsonFieldName := strings.Split(jsonTag, ",")[0]
	return jsonFieldName, true
}

package request

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validates all kind of structs
func Validator(v interface{}) error {
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

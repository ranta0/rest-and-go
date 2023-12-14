package request

import "testing"

func TestGetJSONFieldName(t *testing.T) {
	type TestStruct struct {
		Field1 string `json:"field_1"`
		Field2 int    `json:"field_2"`
		Field3 bool   `json:"field_3"`
	}

	tests := []struct {
		fieldName       string
		expectedJSONTag string
		expectFound     bool
	}{
		{"Field1", "field_1", true},
		{"Field2", "field_2", true},
		{"Field3", "field_3", true},
		{"NonExistentField", "", false},
	}

	for _, test := range tests {
		jsonFieldName, found := getJSONFieldName(TestStruct{}, test.fieldName)

		if found && jsonFieldName != test.expectedJSONTag {
			t.Errorf("Expected JSON field name for %s to be %s, got %s", test.fieldName, test.expectedJSONTag, jsonFieldName)
		}
	}
}

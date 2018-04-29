package main

import (
	"fmt"
)

func readString(data map[string]interface{}, field string) (string, error) {
	value, exists := data[field]
	if !exists {
		return "", fmt.Errorf("Field '%v' is missing", field)
	}

	switch value.(type) {
	case string:
		return value.(string), nil
	default:
		return "", fmt.Errorf("Non-string value for field '%v'", field)
	}
}

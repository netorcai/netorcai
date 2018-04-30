package main

import (
	"fmt"
	"strconv"
)

func readString(data map[string]interface{}, field string) (string, error) {
	value, exists := data[field]
	if !exists {
		return "", fmt.Errorf("Field '%v' is missing", field)
	}

	switch value.(type) {
	default:
		return "", fmt.Errorf("Non-string value for field '%v'", field)
	case string:
		return value.(string), nil
	}
}

func readInt(data map[string]interface{}, field string) (int, error) {
	value, exists := data[field]
	if !exists {
		return 0, fmt.Errorf("Field '%v' is missing", field)
	}

	switch value.(type) {
	default:
		return 0, fmt.Errorf("Non-integral value for field '%v'", field)
	case float64:
		return int(value.(float64)), nil
	}
}

func readObject(data map[string]interface{}, field string) (map[string]interface{}, error) {
	value, exists := data[field]
	if !exists {
		return make(map[string]interface{}),
			fmt.Errorf("Field '%v' is missing", field)
	}

	switch value.(type) {
	default:
		return make(map[string]interface{}),
			fmt.Errorf("Non-object value for field '%v'", field)
	case map[string]interface{}:
		return value.(map[string]interface{}), nil
	}
}

func readIntInString(data map[string]interface{}, field string, bitSize,
	minValue, maxValue int) (int, error) {
	value, exists := data[field]
	if !exists {
		return 0, fmt.Errorf("Field '%v' is missing", field)
	}

	switch value.(type) {
	default:
		return 0, fmt.Errorf("Non-string value for field '%v'", field)
	case string:
		intValue, err := strconv.ParseInt(value.(string), 0, bitSize)
		if err != nil {
			return 0, fmt.Errorf("Field '%v' is invalid: "+
				"Could not parse integer. Err: %v", field, err)
		}

		if intValue < int64(minValue) {
			return int(intValue), fmt.Errorf("Field '%v' is invalid: "+
				"Value is less than minValue=%v",
				field, minValue)
		}

		if intValue > int64(maxValue) {
			return int(intValue), fmt.Errorf("Field '%v' is invalid: "+
				"Value is greater than maxValue=%v",
				field, maxValue)
		}

		return int(intValue), nil

	}
}

func readFloatInString(data map[string]interface{}, field string, bitSize int,
	minValue, maxValue float64) (float64, error) {
	value, exists := data[field]
	if !exists {
		return 0, fmt.Errorf("Field '%v' is missing", field)
	}

	switch value.(type) {
	default:
		return 0, fmt.Errorf("Non-string value for field '%v'", field)
	case string:
		floatValue, err := strconv.ParseFloat(value.(string), bitSize)
		if err != nil {
			return 0, fmt.Errorf("Field '%v' is invalid: "+
				"Could not parse float. Err: %v", field, err)
		}

		if floatValue < minValue {
			return floatValue, fmt.Errorf("Field '%v' is invalid: "+
				"Value is less than minValue=%v",
				field, minValue)
		}

		if floatValue > maxValue {
			return floatValue, fmt.Errorf("Field '%v' is invalid: "+
				"Value is greater than maxValue=%v",
				field, maxValue)
		}

		return floatValue, nil

	}
}

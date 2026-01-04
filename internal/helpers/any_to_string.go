package helpers

import (
	"errors"
	"fmt"
)

func arrayToString[T any](array []T) (string, error) {
	output := ""

	for _, value := range array {
		strVal, err := AnyToString(value)
		if err != nil {
			return "", err
		}

		output += fmt.Sprintf("%s,", strVal)
	}

	output = output[:len(output)-1]
	return output, nil
}

// AnyToString converts any value to a string.
func AnyToString(value any) (string, error) {
	switch value := value.(type) {
	case string:
		return value, nil
	case bool:
		return fmt.Sprintf("%t", value), nil
	case int, int64:
		return fmt.Sprintf("%d", value), nil
	case float32, float64:
		return fmt.Sprintf("%.2f", value), nil
	case []string:
		return arrayToString(value)
	case []bool:
		return arrayToString(value)
	case []int:
		return arrayToString(value)
	case []int64:
		return arrayToString(value)
	case []float32:
		return arrayToString(value)
	case []float64:
		return arrayToString(value)
	default:
		return "", errors.New("undefined type")
	}
}

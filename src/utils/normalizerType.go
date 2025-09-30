package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func ToInt(value interface{}) int {
	if value == nil {
		return 0
	}

	switch v := value.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case int32:
		return int(v)
	case int64:
		return int(v)
	case string:
		v = stringFormat(v)
		if v == "" {
			return 0
		}
		num, err := strconv.Atoi(v)
		if err != nil {
			return 0
		}
		return num
	default:
		return 0
	}
}

func ToFloat64(value interface{}) float64 {
	if value == nil {
		return 0.0
	}

	switch v := value.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case string:
		v = stringFormat(v)
		if v == "" {
			return 0
		}
		num, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0
		}
		return num
	default:
		return 0
	}
}

func ToString(value interface{}) string {
	if value == nil {
		return ""
	}
	return stringFormat(fmt.Sprintf("%v", value))
}

func stringFormat(datatype string) string {
	return strings.TrimSpace(strings.ReplaceAll(datatype, ",", ""))
}

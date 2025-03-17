package provider

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// parseQuery extracts a query parameter of type T with a default fallback and error handling
func ParseQuery[T any](r *http.Request, key string, defaultValue T) (T, error) {
	values := r.URL.Query()

	if val, exists := values[key]; exists && len(val) > 0 {
		var result T
		switch any(result).(type) {
		case int:
			parsed, err := strconv.Atoi(val[0])
			if err != nil {
				log.Printf("error parsing int for key '%s': %v", key, err)
				return defaultValue, err
			}
			return any(parsed).(T), nil
		case float64:
			parsed, err := strconv.ParseFloat(val[0], 64)
			if err != nil {
				log.Printf("error parsing float for key '%s': %v", key, err)
				return defaultValue, err
			}
			return any(parsed).(T), nil
		case bool:
			parsed, err := strconv.ParseBool(val[0])
			if err != nil {
				log.Printf("error parsing bool for key '%s': %v", key, err)
				return defaultValue, err
			}
			return any(parsed).(T), nil
		case string:
			return any(val[0]).(T), nil
		default:
			log.Printf("unsupported type for key '%s'", key)
			return defaultValue, errors.New("unsupported type")
		}
	}
	return defaultValue, nil
}

// parseQueryList extracts a list of type T from a comma-separated query parameter with error handling
func ParseQueryList[T any](r *http.Request, key, separator string) ([]T, error) {
	values := r.URL.Query()

	if val, exists := values[key]; exists && len(val) > 0 {
		strList := strings.Split(val[0], separator)
		var result []T
		var parseErr error

		for _, strVal := range strList {
			var converted T
			switch any(converted).(type) {
			case int:
				parsed, err := strconv.Atoi(strings.TrimSpace(strVal))
				if err != nil {
					log.Printf("Error parsing int list for key '%s': %v", key, err)
					parseErr = err
					continue
				}
				converted = any(parsed).(T)
			case float64:
				parsed, err := strconv.ParseFloat(strings.TrimSpace(strVal), 64)
				if err != nil {
					log.Printf("Error parsing float list for key '%s': %v", key, err)
					parseErr = err
					continue
				}
				converted = any(parsed).(T)
			case string:
				converted = any(strings.TrimSpace(strVal)).(T)
			default:
				log.Printf("Unsupported type for key '%s'", key)
				parseErr = errors.New("unsupported type")
				continue
			}
			result = append(result, converted)
		}
		if parseErr != nil {
			return result, parseErr
		}
		return result, nil
	}
	return nil, errors.New("missing or empty query parameter")
}

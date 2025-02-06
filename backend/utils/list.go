package utils

import "errors"

func FirstWhere[T any](slice []T, condition func(T) bool) (T, error) {
	for _, item := range slice {
		if condition(item) {
			return item, nil
		}
	}
	var zero T // Zero value of type T
	return zero, errors.New("no matching element found")
}

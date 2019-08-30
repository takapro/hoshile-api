package main

import (
	"errors"
	"regexp"
	"strconv"
)

var (
	ErrInvalidUrlParam = errors.New("invalid url param")
)

func ParseUrlParam(path string) (int, error) {
	result := regexp.MustCompile(`^/\w+/(\d+)$`).FindStringSubmatch(path)
	if result == nil || len(result) != 2 {
		return 0, ErrInvalidUrlParam
	}

	return strconv.Atoi(result[1])
}

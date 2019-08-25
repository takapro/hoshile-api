package main

import (
	"encoding/json"
	"errors"
	"net/http"
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

func WriteJson(w http.ResponseWriter, obj interface{}) {
	b, err := json.Marshal(obj)
	if err != nil {
		http.Error(w, "JSON error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

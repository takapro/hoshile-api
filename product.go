package main

import (
	"net/http"
)

type Product struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Brand    string  `json:"brand"`
	Price    float64 `json:"price"`
	ImageUrl string  `json:"imageUrl"`
}

func HandleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getProducts(w)

	case http.MethodOptions:
		HandleOptions(w)

	default:
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
}

func HandleProduct(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getProduct(w, r)

	case http.MethodOptions:
		HandleOptions(w)

	default:
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
}

func getProducts(w http.ResponseWriter) {
	arr, err := SelectProducts()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	WriteJson(w, arr)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	id, err := ParseUrlParam(r.URL.Path)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	p, err := SelectProduct(id)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	WriteJson(w, p)
}

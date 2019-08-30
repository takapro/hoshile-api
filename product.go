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

func GetProducts(r *http.Request) (interface{}, error) {
	return SelectProducts()
}

func GetProduct(r *http.Request) (interface{}, error) {
	id, err := ParseUrlParam(r.URL.Path)
	if err != nil {
		return nil, ErrBadRequest
	}

	return SelectProduct(id)
}

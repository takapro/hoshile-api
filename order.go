package main

import (
	"encoding/json"
	"net/http"
)

type OrderHead struct {
	Id         int           `json:"id"`
	UserId     int           `json:"userId"`
	CreateDate string        `json:"createDate"`
	Details    []OrderDetail `json:"details"`
}

type OrderDetail struct {
	Id        int      `json:"id"`
	OrderId   int      `json:"orderId"`
	ProductId int      `json:"productId"`
	Quantity  int      `json:"quantity"`
	Product   *Product `json:"product"`
}

type OrderParams struct {
	ProductId int `json:"productId"`
	Quantity  int `json:"quantity"`
}

type orderIdParams struct {
	Id int `json:"id"`
}

func GetOrders(w http.ResponseWriter, r *http.Request) {
	userId, _, ok := FindSession(r)
	if !ok {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	orders, err := SelectOrderHeads(userId)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	for _, order := range orders {
		details, err := getDetails(order.Id)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		order.Details = details
	}

	WriteJson(w, orders)
}

func GetOrder(w http.ResponseWriter, r *http.Request) {
	userId, _, ok := FindSession(r)
	if !ok {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var p orderIdParams
	if r.Body == nil || json.NewDecoder(r.Body).Decode(&p) != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	order, err := SelectOrderHead(p.Id)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if order.UserId != userId {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	details, err := getDetails(order.Id)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	order.Details = details

	WriteJson(w, order)
}

func PostOrder(w http.ResponseWriter, r *http.Request) {
	userId, _, ok := FindSession(r)
	if !ok {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var params []OrderParams
	if r.Body == nil || json.NewDecoder(r.Body).Decode(&params) != nil || len(params) == 0 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	orderId, err := InsertOrder(userId, params)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	WriteJson(w, orderId)
}

func getDetails(orderId int) ([]OrderDetail, error) {
	details, err := SelectOrderDetails(orderId)
	if err != nil {
		return nil, err
	}

	for _, detail := range details {
		product, err := SelectProduct(detail.ProductId)
		if err != nil {
			return nil, err
		}

		detail.Product = product
	}

	return details, nil
}

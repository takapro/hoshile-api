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

func GetOrders(r *http.Request) (interface{}, error) {
	userId, _, ok := FindSession(r)
	if !ok {
		return nil, ErrForbidden
	}

	orders, err := SelectOrderHeads(userId)
	if err != nil {
		return nil, err
	}

	for i := range orders {
		details, err := getDetails(orders[i].Id)
		if err != nil {
			return nil, err
		}

		orders[i].Details = details
	}

	return orders, nil
}

func GetOrder(r *http.Request) (interface{}, error) {
	userId, _, ok := FindSession(r)
	if !ok {
		return nil, ErrForbidden
	}

	orderId, err := ParseUrlParam(r.URL.Path)
	if err != nil {
		return nil, ErrBadRequest
	}

	order, err := SelectOrderHead(orderId)
	if err != nil {
		return nil, err
	}
	if order.UserId != userId {
		return nil, ErrBadRequest
	}

	details, err := getDetails(orderId)
	if err != nil {
		return nil, err
	}

	order.Details = details

	return order, nil
}

func PostOrder(r *http.Request) (interface{}, error) {
	userId, _, ok := FindSession(r)
	if !ok {
		return nil, ErrForbidden
	}

	var params []OrderParams
	if r.Body == nil || json.NewDecoder(r.Body).Decode(&params) != nil || len(params) == 0 {
		return nil, ErrBadRequest
	}

	return InsertOrder(userId, params)
}

func getDetails(orderId int) ([]OrderDetail, error) {
	details, err := SelectOrderDetails(orderId)
	if err != nil {
		return nil, err
	}

	for i := range details {
		product, err := SelectProduct(details[i].ProductId)
		if err != nil {
			return nil, err
		}

		details[i].Product = product
	}

	return details, nil
}

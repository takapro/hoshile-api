package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Token        string `json:"token"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	ShoppingCart string `json:"shoppingCart"`
	IsAdmin      bool   `json:"isAdmin"`
}

type userParams struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type passwordParams struct {
	CurPassword string `json:"curPassword"`
	NewPassword string `json:"newPassword"`
}

type shoppingCartParams struct {
	ShoppingCart string `json:"shoppingCart"`
}

func UserLogin(r *http.Request) (interface{}, error) {
	var p userParams
	if r.Body == nil || json.NewDecoder(r.Body).Decode(&p) != nil || p.Email == "" || p.Password == "" {
		return nil, ErrBadRequest
	}

	user, id, err := AuthUser(0, p.Email, p.Password)
	if err != nil {
		return nil, err
	}

	user.Token = NewSession(id)

	return user, nil
}

func UserSignup(r *http.Request) (interface{}, error) {
	var p userParams
	if r.Body == nil || json.NewDecoder(r.Body).Decode(&p) != nil || p.Name == "" || p.Email == "" || p.Password == "" {
		return nil, ErrBadRequest
	}

	id, err := InsertUser(p.Name, p.Email, p.Password)
	if err != nil {
		return nil, err
	}

	return getUser(id, NewSession(id))
}

func GetProfile(r *http.Request) (interface{}, error) {
	id, token, ok := FindSession(r)
	if !ok {
		return nil, ErrBadRequest
	}

	return getUser(id, token)
}

func PutProfile(r *http.Request) (interface{}, error) {
	id, token, ok := FindSession(r)
	if !ok {
		return nil, ErrBadRequest
	}

	var p userParams
	if r.Body == nil || json.NewDecoder(r.Body).Decode(&p) != nil || p.Name == "" || p.Email == "" {
		return nil, ErrBadRequest
	}

	err := UpdateUser(id, map[string]string{"name": p.Name, "email": p.Email})
	if err != nil {
		return nil, err
	}

	return getUser(id, token)
}

func PutPassword(r *http.Request) (interface{}, error) {
	id, token, ok := FindSession(r)
	if !ok {
		return nil, ErrBadRequest
	}

	var p passwordParams
	if r.Body == nil || json.NewDecoder(r.Body).Decode(&p) != nil || p.CurPassword == "" || p.NewPassword == "" {
		return nil, ErrBadRequest
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(p.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, id, err := AuthUser(id, "", p.CurPassword)
	if err != nil {
		return nil, err
	}

	err = UpdateUser(id, map[string]string{"password": string(hash)})
	if err != nil {
		return nil, err
	}

	user.Token = token

	return user, nil
}

func PutShoppingCart(r *http.Request) (interface{}, error) {
	id, token, ok := FindSession(r)
	if !ok {
		return nil, ErrBadRequest
	}

	var p shoppingCartParams
	if r.Body == nil || json.NewDecoder(r.Body).Decode(&p) != nil {
		return nil, ErrBadRequest
	}

	err := UpdateUser(id, map[string]string{"shoppingCart": p.ShoppingCart})
	if err != nil {
		return nil, err
	}

	return getUser(id, token)
}

func getUser(id int, token string) (interface{}, error) {
	user, err := SelectUser(id)
	if err != nil {
		return nil, err
	}

	user.Token = token

	return user, nil
}

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

func UserLogin(w http.ResponseWriter, r *http.Request) {
	var p userParams
	if r.Body == nil || json.NewDecoder(r.Body).Decode(&p) != nil || p.Email == "" || p.Password == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	user, id, err := AuthUser(0, p.Email, p.Password)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	user.Token = NewSession(id)

	WriteJson(w, user)
}

func UserSignup(w http.ResponseWriter, r *http.Request) {
	var p userParams
	if r.Body == nil || json.NewDecoder(r.Body).Decode(&p) != nil || p.Name == "" || p.Email == "" || p.Password == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	id, err := InsertUser(p.Name, p.Email, p.Password)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	writeUser(w, id, NewSession(id))
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	id, token, ok := FindSession(r)
	if !ok {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	writeUser(w, id, token)
}

func PutProfile(w http.ResponseWriter, r *http.Request) {
	id, token, ok := FindSession(r)
	if !ok {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var p userParams
	if r.Body == nil || json.NewDecoder(r.Body).Decode(&p) != nil || p.Name == "" || p.Email == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err := UpdateUser(id, map[string]string{"name": p.Name, "email": p.Email})
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	writeUser(w, id, token)
}

func PutPassword(w http.ResponseWriter, r *http.Request) {
	id, token, ok := FindSession(r)
	if !ok {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var p passwordParams
	if r.Body == nil || json.NewDecoder(r.Body).Decode(&p) != nil || p.CurPassword == "" || p.NewPassword == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(p.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	user, id, err := AuthUser(id, "", p.CurPassword)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	err = UpdateUser(id, map[string]string{"password": string(hash)})
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	user.Token = token

	WriteJson(w, user)
}

func PutShoppingCart(w http.ResponseWriter, r *http.Request) {
	id, token, ok := FindSession(r)
	if !ok {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var p shoppingCartParams
	if r.Body == nil || json.NewDecoder(r.Body).Decode(&p) != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err := UpdateUser(id, map[string]string{"shoppingCart": p.ShoppingCart})
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	writeUser(w, id, token)
}

func writeUser(w http.ResponseWriter, id int, token string) {
	user, err := SelectUser(id)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	user.Token = token

	WriteJson(w, user)
}

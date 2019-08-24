package main

import (
	"encoding/json"
	"net/http"
)

type User struct {
	Session      string `json:"session"`
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

	id, user, err := AuthUser(p.Email, p.Password)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	user.Session = NewSession(id)

	WriteJson(w, user)
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	id, ok := FindSession(r)
	if !ok {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	user, err := SelectUser(id)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	WriteJson(w, user)
}

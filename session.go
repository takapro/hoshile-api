package main

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"
)

var c *Cache

func InitSession() {
	c = NewCache(24 * 60 * 60)
}

func NewSession(id int) string {
	buf := make([]byte, 32)
	_, err := rand.Read(buf)
	if err != nil {
		return ""
	}

	token := strings.TrimRight(base64.URLEncoding.EncodeToString(buf), "=")
	c.Set(token, id)
	return token
}

func FindSession(r *http.Request) (int, string, bool) {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return 0, "", false
	}

	token := strings.TrimPrefix(auth, "Bearer ")
	id, ok := c.Get(token)
	if !ok {
		return 0, "", false
	}

	return id, token, true
}

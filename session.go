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

	session := strings.TrimRight(base64.URLEncoding.EncodeToString(buf), "=")
	c.Set(session, id)
	return session
}

func FindSession(r *http.Request) (int, bool) {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return 0, false
	}

	session := strings.TrimPrefix(auth, "Bearer ")
	return c.Get(session)
}

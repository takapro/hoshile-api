package main

import (
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	err := InitDB()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", hello)

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

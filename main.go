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

	http.HandleFunc("/products", HandleProducts)
	http.HandleFunc("/products/", HandleProduct)

	http.HandleFunc("/", handleNotFound)

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not found", http.StatusNotFound)
}

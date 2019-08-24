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

	router := NewRouter()

	router.Handle("/products", Handlers{http.MethodGet: GetProducts})
	router.Handle("/products/", Handlers{http.MethodGet: GetProduct})

	router.Handle("/user/login", Handlers{http.MethodPost: UserLogin})

	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		panic(err)
	}
}

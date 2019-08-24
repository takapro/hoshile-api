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

	InitSession()

	router := NewRouter()

	router.Handle("/products", Handlers{http.MethodGet: GetProducts})
	router.Handle("/products/", Handlers{http.MethodGet: GetProduct})

	router.Handle("/user/login", Handlers{http.MethodPost: UserLogin})
	router.Handle("/user/profile", Handlers{http.MethodGet: GetProfile})

	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		panic(err)
	}
}

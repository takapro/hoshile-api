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
	router.Handle("/user/signup", Handlers{http.MethodPost: UserSignup})
	router.Handle("/user/profile", Handlers{http.MethodGet: GetProfile, http.MethodPut: PutProfile})
	router.Handle("/user/password", Handlers{http.MethodPut: PutPassword})
	router.Handle("/user/shoppingCart", Handlers{http.MethodPut: PutShoppingCart})

	router.Handle("/orders", Handlers{http.MethodGet: GetOrders, http.MethodPost: PostOrder})
	router.Handle("/orders/", Handlers{http.MethodGet: GetOrder})

	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		panic(err)
	}
}

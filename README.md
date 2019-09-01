# HoshiLe’s Store API Server in Go

`hoshile-api` is a Go port of the API server portion of HoshiLe’s Store.

## About HoshiLe’s Store

HoshiLe’s Store is a classroom project written by [Ngoc Tin Le](https://github.com/takint) and [Takanori Hoshi](https://github.com/takapro).
Original was written in PHP and composed of a Rest API server, store front and admin clients.

## API endpoints

- `GET /products`

- `GET /products/:id`

- `POST /user/login`

- `POST /user/signup`

- `GET /user/profile` (requires login)

- `PUT /user/profile` (requires login)

- `PUT /user/password` (requires login)

- `PUT /user/shoppingCart` (requires login)

- `GET /orders` (requires login)

- `GET /orders/:id` (requires login)

- `POST /orders` (requires login)

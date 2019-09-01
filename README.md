# HoshiLe’s Store API Server in Go

`hoshile-api` is a Go port of the API server portion of HoshiLe’s Store.

## About HoshiLe’s Store

HoshiLe’s Store is a classroom project written by [Ngoc Tin Le](https://github.com/takint) and [Takanori Hoshi](https://github.com/takapro).
Original was written in PHP and composed of a Rest API server, store front and admin clients.

## API endpoints

- `GET /products`

    Returns the list of all products in the database.

    ```
    [
      {
        "id": 1,
        "name": "MacBook Air",
        "brand": "Apple",
        "price": 1234.56,
        "imageUrl": "https://example.com/images/product-1.jpg"
      },
      ...
    ]
    ```

- `GET /products/:id`

    Returns a single product with the specified product id.

    ```
    {
      "id": 1,
      "name": "MacBook Air",
      "brand": "Apple",
      "price": 1234.56,
      "imageUrl": "https://example.com/images/product-1.jpg"
    }
    ```

- `POST /user/login`

    Requires the parameters `email` and `password`.

    ```
    {
      "email": "taka@email.ca",
      "password": "1234"
    }
    ```

    Returns a user object upon a successful login.

    ```
    {
      "token": "SECRET_TOKEN",
      "name": "Takanori Hoshi",
      "email": "taka@email.ca",
      "shoppingCart": "[{\"productId\":1,\"quantity\":2}]",
      "isAdmin": false
    }
    ```

    After login, each endpoint which requires login can be accessed with the authorization header.

    ```
    Authorization: Bearer SECRET_TOKEN
    ```

- `POST /user/signup`

    Requires the parameters `name`, `email` and `password`.

    ```
    {
      "name": "Takanori Hoshi",
      "email": "taka@email.ca",
      "password": "1234"
    }
    ```

    Returns a newly created user object (same as POST /login).

- `GET /user/profile` (requires login)

    Returns the user object of the user who is logged in.

- `PUT /user/profile` (requires login)

    Update the user's profile (`name` and `email`).

    ```
    {
      "name": "Takanori Hoshi",
      "email": "taka@email.ca"
    }
    ```

    Returns the updated user object.

- `PUT /user/password` (requires login)

    Update the user's password.

    ```
    {
      "curPassword": "1234",
      "newPassword": "5678"
    }
    ```

    Returns the updated user object.

- `PUT /user/shoppingCart` (requires login)

    Update the user's shopping cart.

    ```
    {
      "shoppingCart": "[{\"productId\":1,\"quantity\":2},{\"productId\":3,\"quantity\":1}]"
    }
    ```

    Returns `true` on success.

    ```
    true
    ```

- `GET /orders` (requires login)

    Returns the list of all orders of the user.

    ```
    [
      {
        "id": 3,
        "userId": 4,
        "createDate": "2019-09-01 12:34:56",
        "details": [
          {
            "id": 5,
            "orderId": 3,
            "productId": 1,
            "quantity": 2,
            "product": {
              "id": 1,
              "name": "MacBook Air",
              "brand": "Apple",
              "price": 1234.56,
              "imageUrl": "https://example.com/images/product-1.jpg"
            }
          },
          ...
        ]
      },
      ...
    ]
    ```

- `GET /orders/:id` (requires login)

    Returns a single order with the specified order id.

    ```
    {
      "id": 3,
      "userId": 4,
      "createDate": "2019-09-01 12:34:56",
      "details": [
        {
          "id": 5,
          "orderId": 3,
          "productId": 1,
          "quantity": 2,
          "product": {
            "id": 1,
            "name": "MacBook Air",
            "brand": "Apple",
            "price": 1234.56,
            "imageUrl": "https://example.com/images/product-1.jpg"
          }
        },
        ...
      ]
    }
    ```

- `POST /orders` (requires login)

    Requires an array of objects with `productId` and `quantity`.

    ```
    [
      {
        "productId": 1,
        "quantity": 2
      },
      {
        "productId": 3,
        "quantity": 1
      }
    ]
    ```

    Returns the id of the newly created order object.

    ```
    3
    ```

## How to run

Build with `go build`.

```
$ go build
```

Create a database on MySQL.

```
$ mysql -u root -e "create database HoshiLe"
$ mysql -u root HoshiLe < sql/tables.sql
$ mysql -u root HoshiLe < sql/inserts_data.sql
```

Just run (default port is 3000).

```
$ ./hoshile-api
```

Run in debug mode to display each response, and with 500ms of delay for each response.

```
$ ./hoshile-api --debug --delay 500
```

Specify the API port and the database parameters.

```
$ export PORT=3456
$ export DB_HOST="db.example.com"
$ export DB_USER="myuser"
$ export DB_PASS="password"
$ export DB_NAME="mydb"
$ ./hoshile-api
```

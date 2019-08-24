package main

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func InitDB() error {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	name := os.Getenv("DB_NAME")
	if host == "" {
		host = "localhost"
	}
	if user == "" {
		user = "root"
	}
	if name == "" {
		name = "HoshiLe"
	}
	dsn := user + ":" + pass + "@tcp(" + host + ")/" + name

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	err = db.Ping()
	return err
}

func SelectProducts() ([]Product, error) {
	rows, err := db.Query("select id, name, brand, price, imageUrl from Products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	arr := []Product{}
	for rows.Next() {
		var p Product
		err = rows.Scan(&p.Id, &p.Name, &p.Brand, &p.Price, &p.ImageUrl)
		if err != nil {
			return nil, err
		}
		arr = append(arr, p)
	}

	return arr, nil
}

func SelectProduct(id int) (*Product, error) {
	row := db.QueryRow("select id, name, brand, price, imageUrl from Products where id = ?", id)

	var p Product
	err := row.Scan(&p.Id, &p.Name, &p.Brand, &p.Price, &p.ImageUrl)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func SelectUser(id int) (*User, error) {
	row := db.QueryRow("select name, email, shoppingCart, isAdmin = b'1' from Users where id = ?", id)

	var u User
	err := row.Scan(&u.Name, &u.Email, &u.ShoppingCart, &u.IsAdmin)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func AuthUser(email string, password string) (int, *User, error) {
	row := db.QueryRow("select id, name, email, password, shoppingCart, isAdmin = b'1' from Users where email = ?", email)

	var u User
	var id int
	var hash string
	err := row.Scan(&id, &u.Name, &u.Email, &hash, &u.ShoppingCart, &u.IsAdmin)
	if err != nil {
		return 0, nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return 0, nil, err
	}

	return id, &u, nil
}

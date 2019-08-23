package main

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
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
	rows, err := db.Query("select * from Products")
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
	row := db.QueryRow("select * from Products where id = ?", id)

	var p Product
	err := row.Scan(&p.Id, &p.Name, &p.Brand, &p.Price, &p.ImageUrl)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

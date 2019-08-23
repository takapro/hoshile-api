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

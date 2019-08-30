package main

import (
	"database/sql"
	"os"
	"strings"

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
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func SelectUser(id int) (*User, error) {
	row := db.QueryRow("select name, email, shoppingCart, isAdmin = b'1' from Users where id = ?", id)

	var u User
	var shoppingCart sql.NullString
	err := row.Scan(&u.Name, &u.Email, &shoppingCart, &u.IsAdmin)
	if err == sql.ErrNoRows {
		return nil, ErrBadRequest
	}
	if err != nil {
		return nil, err
	}

	if shoppingCart.Valid {
		u.ShoppingCart = shoppingCart.String
	}

	return &u, nil
}

func AuthUser(id int, email string, password string) (*User, int, error) {
	query := "select id, name, email, password, shoppingCart, isAdmin = b'1' from Users where "
	var row *sql.Row
	if email == "" {
		row = db.QueryRow(query+"id = ?", id)
	} else {
		row = db.QueryRow(query+"email = ?", email)
	}

	var u User
	var hash string
	var shoppingCart sql.NullString
	err := row.Scan(&id, &u.Name, &u.Email, &hash, &shoppingCart, &u.IsAdmin)
	if err == sql.ErrNoRows {
		return nil, 0, ErrBadRequest
	}
	if err != nil {
		return nil, 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return nil, 0, ErrBadRequest
	}

	if shoppingCart.Valid {
		u.ShoppingCart = shoppingCart.String
	}

	return &u, id, nil
}

func InsertUser(name, email, password string) (int, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	result, err := db.Exec("insert into Users (name, email, password) values (?, ?, ?)", name, email, hash)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return int(id), err
}

func UpdateUser(id int, kv map[string]string) error {
	keys := []string{}
	values := []interface{}{}
	for k, v := range kv {
		keys = append(keys, k)
		values = append(values, v)
	}
	values = append(values, id)

	query := "update Users set " + strings.Join(keys, " = ?, ") + " = ? where id = ?"
	_, err := db.Exec(query, values...)
	return err
}

func SelectOrderHeads(userId int) ([]OrderHead, error) {
	rows, err := db.Query("select id, userId, createDate from OrderHeads where userId = ?", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	arr := []OrderHead{}
	for rows.Next() {
		var o OrderHead
		err = rows.Scan(&o.Id, &o.UserId, &o.CreateDate)
		if err != nil {
			return nil, err
		}

		arr = append(arr, o)
	}

	return arr, nil
}

func SelectOrderHead(id int) (*OrderHead, error) {
	row := db.QueryRow("select id, userId, createDate from OrderHeads where id = ?", id)

	var o OrderHead
	err := row.Scan(&o.Id, &o.UserId, &o.CreateDate)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &o, nil
}

func SelectOrderDetails(orderId int) ([]OrderDetail, error) {
	rows, err := db.Query("select id, orderId, productId, quantity from OrderDetails where orderId = ?", orderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	arr := []OrderDetail{}
	for rows.Next() {
		var d OrderDetail
		err = rows.Scan(&d.Id, &d.OrderId, &d.ProductId, &d.Quantity)
		if err != nil {
			return nil, err
		}

		arr = append(arr, d)
	}

	return arr, nil
}

func InsertOrder(userId int, params []OrderParams) (int, error) {
	result, err := db.Exec("insert into OrderHeads (userId, createDate) values (?, now())", userId)
	if err != nil {
		return 0, err
	}

	orderId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	for _, p := range params {
		_, err := db.Exec("insert into OrderDetails (orderId, productId, quantity) values (?, ?, ?)", orderId, p.ProductId, p.Quantity)
		if err != nil {
			return 0, err
		}
	}

	return int(orderId), nil
}

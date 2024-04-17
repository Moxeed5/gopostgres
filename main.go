package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Product struct {
	Name      string
	Price     float64
	Available bool
}

func main() {
	connStr := "postgres://postgres:secret@localhost:5432/gopgtest?sslmode=disable"

	db, err := sql.Open("postgres", connStr)

	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	if db.Ping(); err != nil {
		log.Fatal(err)
	}

	createProductTable(db)

	data := []Product{}
	rows, err := db.Query("SELECT name, available, price FROM product")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var name string
	var available bool
	var price float64

	for rows.Next() {
		err := rows.Scan(&name, &available, &price)
		if err != nil {
			log.Fatal(err)
		}

		data = append(data, Product{name, price, available})
	}

	fmt.Println(data)

}

func createProductTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS product (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		price NUMERIC(6,2) NOT NULL,
		created timestamp DEFAULT NOW()	

	)`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func insertProduct(db *sql.DB, product Product) int {
	//returning int bc we want to get primary key of product we inserted

	query := `INSERT INTO product (name,price, available)
		VALUES ($1,$2, $3) RETURNING id`

	var pk int //store primary key that we're returning in the func
	err := db.QueryRow(query, product.Name, product.Price, product.Available).Scan(&pk)
	if err != nil {
		log.Fatal(err)
	}

	return pk

}

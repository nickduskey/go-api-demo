package models

import (
	"database/sql"
	"fmt"
	"log"
)

// Product represents a product
type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// GetProduct fetches a product from the db
func (p *Product) GetProduct(db *sql.DB) error {
	log.Printf("GetProducts() called")
	return db.QueryRow("SELECT name, price FROM products WHERE id=$1",
		p.ID).Scan(&p.Name, &p.Price)
}

// UpdateProduct updates a product by id in the db
func (p *Product) UpdateProduct(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE products SET name=$1, price=$2 where id=$3",
			p.Name, p.Price, p.ID)

	return err
}

// DeleteProduct deletes a product
func (p *Product) DeleteProduct(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM products WHERE id=$1", p.ID)

	return err
}

// CreateProduct creates a product
func (p *Product) CreateProduct(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO products(name, price) VALUES($1, $2) RETURNING id",
		p.Name, p.Price).Scan(&p.ID)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// GetProducts fetches products from the db
func GetProducts(db *sql.DB, start, count int) ([]Product, error) {
	rows, err := db.Query(
		"SELECT id, name, price FROM products LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	products := []Product{}

	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

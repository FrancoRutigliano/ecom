package product

import (
	"database/sql"
	"fmt"
	"strings"

	typePayload "github.com/FrancoRutigliano/ecom/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetProducts() ([]typePayload.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}

	products := make([]typePayload.Product, 0)
	for rows.Next() {
		p, err := scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, *p)
	}
	return products, nil
}

func (s *Store) GetProductsByIDs(productIDs []int) ([]typePayload.Product, error) {
	placeholder := strings.Repeat(",?", len(productIDs)-1)
	query := fmt.Sprintf("SELECT * FROM products WHERE id in (?%s)", placeholder)

	//convert productIDs into []interfaces{}
	args := make([]interface{}, len(productIDs))
	for i, v := range productIDs {
		args[i] = v
	}

	// query to database passing the query and the args...
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	// making a []typepayload.Product{}
	products := []typePayload.Product{}
	for rows.Next() {
		p, err := scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, *p)
	}

	return products, nil
}

func (s *Store) CreateProduct(product typePayload.Product) error {
	_, err := s.db.Exec("INSERT INTO products (name, description, image, price, quantity) VALUES (?, ?, ?, ?, ?)", product.Name, product.Description, product.Image, product.Price, product.Quantity)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) UpdateProduct(product typePayload.Product) error {
	_, err := s.db.Exec("UPDATE products SET name = ?, price = ?, image = ?, description = ?, quantity = ? WHERE id = ?",
		product.Name, product.Price, product.Image, product.Description, product.Quantity, product.ID)
	if err != nil {
		return err
	}
	return nil
}

func scanRowsIntoProduct(rows *sql.Rows) (*typePayload.Product, error) {
	product := new(typePayload.Product)

	err := rows.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Image,
		&product.Price,
		&product.Quantity,
		&product.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return product, nil

}

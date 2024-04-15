// model.go

package main

import (
	"database/sql"
)

type product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (p *product) getProduct(db *sql.DB) error {
	return db.QueryRow("SELECT name, price FROM products WHERE id=$1",
		p.ID).Scan(&p.Name, &p.Price)
}

func (p *product) updateProduct(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE products SET name=$1, price=$2 WHERE id=$3",
			p.Name, p.Price, p.ID)

	return err
}

func (p *product) deleteProduct(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM products WHERE id=$1", p.ID)

	return err
}

func (p *product) createProduct(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO products(name, price) VALUES($1, $2) RETURNING id",
		p.Name, p.Price).Scan(&p.ID)

	if err != nil {
		return err
	}

	return nil
}

func getProducts(db *sql.DB, start, count int) ([]product, error) {
	rows, err := db.Query(
		"SELECT id, name,  price FROM products LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	products := []product{}

	for rows.Next() {
		var p product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

//###########################################################

type tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (t *tag) getTag(db *sql.DB) error {
	return db.QueryRow("SELECT name FROM tag WHERE id=$1",
		t.ID).Scan(&t.Name)
}

// Additional way to retrieve categories, plus check to avoid duplicate names (setting the column to unqiue in the actual database is also possible)
func (t *tag) getTagByName(db *sql.DB) error {
	return db.QueryRow("SELECT id FROM tag WHERE LOWER(name)=LOWER($1)",
		t.Name).Scan(&t.ID)
}

func (t *tag) updateTag(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE tag SET name=$1 WHERE id=$2",
			t.Name, t.ID)

	return err
}

func (t *tag) deleteTag(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM tag WHERE id=$1", t.ID)

	return err
}

func (c *tag) createTag(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO tag(name) VALUES($1) RETURNING id",
		c.Name).Scan(&c.ID)

	if err != nil {
		return err
	}

	return nil
}

func getTags(db *sql.DB, start, count int) ([]tag, error) {
	rows, err := db.Query(
		"SELECT id, name FROM tag LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tags := []tag{}

	for rows.Next() {
		var t tag
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}

	return tags, nil
}

type productToTagAssignment struct {
	ID        int `json:"id"`
	ProductID int `json:"productID"`
	TagID     int `json:"TagID"`
}

func (pta *productToTagAssignment) getProductToTagAssignment(db *sql.DB) error {
	return db.QueryRow("SELECT name FROM productToTagAssignment WHERE id=$1",
		pta.ID).Scan(&pta.ID)
}

func (pta *productToTagAssignment) updateProductToTagAssignment(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE productToTagAssignment SET productID=$1,TagID=$2 WHERE id=$3",
			pta.ProductID, pta.TagID, pta.ID)

	return err
}

func (pta *productToTagAssignment) deleteProductToTagAssignment(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM productToTagAssignment WHERE id=$1", pta.ID)

	return err
}

func (pta *productToTagAssignment) deleteProductToTagAssignmentByProduct(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM productToTagAssignment WHERE productID=$1", pta.ProductID)

	return err
}

func (pta *productToTagAssignment) deleteProductToTagAssignmentByTag(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM productToTagAssignment WHERE TagID=$1", pta.TagID)

	return err
}

func (pta *productToTagAssignment) createProductToTagAssignment(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO productToTagAssignment(productID,TagID) VALUES($1,$2) RETURNING id",
		pta.ProductID, pta.TagID).Scan(&pta.ID)

	if err != nil {
		return err
	}

	return nil
}

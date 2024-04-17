// main_test.go

package main

import (
	"log"
	"os"
	"testing"

	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
)

var a App

func TestMain(m *testing.M) {
	a.Initialize(
		"postgres",
		"postgres",
		"postgres")

	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM products")
	a.DB.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")
	a.DB.Exec("DELETE FROM tag")
	a.DB.Exec("ALTER SEQUENCE tag_id_seq RESTART WITH 1")
	a.DB.Exec("DELETE FROM productToTagAssignment")
	a.DB.Exec("ALTER SEQUENCE productToTagAssignment_id_seq RESTART WITH 1")
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS products
(
    id SERIAL,
    name TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
    CONSTRAINT products_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS tag
(
    id SERIAL,
    name TEXT NOT NULL,
    CONSTRAINT tag_pkey PRIMARY KEY (id)
); 

CREATE TABLE IF NOT EXISTS productToTagAssignment
(
    id SERIAL,
    productID integer NOT NULL,
	tagID integer NOT NULL,
    CONSTRAINT productToTagAssignment_pkey PRIMARY KEY (id),
	CONSTRAINT product_fkey FOREIGN KEY (productID) REFERENCES products (id) ON DELETE CASCADE,
	CONSTRAINT tag_fkey FOREIGN KEY (tagID) REFERENCES tag (id) ON DELETE CASCADE
); 

`

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/products", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n and ", expected, actual)
	}
}

func TestGetNonExistentProduct(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/product/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Product not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Product not found'. Got '%s'", m["error"])
	}
}

func TestCreateProduct(t *testing.T) {

	clearTable()

	var jsonStr = []byte(`{"name":"test product", "price": 11.22}`)
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "test product" {
		t.Errorf("Expected product name to be 'test product'. Got '%v'", m["name"])
	}

	if m["price"] != 11.22 {
		t.Errorf("Expected product price to be '11.22'. Got '%v'", m["price"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected product ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetProduct(t *testing.T) {
	clearTable()
	addProducts(1)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

// main_test.go

func addProducts(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO products(name, price) VALUES($1, $2)", "Product "+strconv.Itoa(i), (i+1.0)*10)
	}
}

func TestUpdateProduct(t *testing.T) {

	clearTable()
	addProducts(1)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req)
	var originalProduct map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalProduct)

	var jsonStr = []byte(`{"name":"test product - updated name", "price": 11.22}`)
	req, _ = http.NewRequest("PUT", "/product/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalProduct["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalProduct["id"], m["id"])
	}

	if m["name"] == originalProduct["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalProduct["name"], m["name"], m["name"])
	}

	if m["price"] == originalProduct["price"] {
		t.Errorf("Expected the price to change from '%v' to '%v'. Got '%v'", originalProduct["price"], m["price"], m["price"])
	}
}

func TestDeleteProduct(t *testing.T) {
	clearTable()
	addProducts(1)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/product/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/product/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

// ##########################################################
func TestGetNonExistentTag(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/tag/999", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Tag not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Tag not found'. Got '%s'", m["error"])
	}
}

func TestCreateTag(t *testing.T) {

	clearTable()

	var jsonStr = []byte(`{"name":"test tag"}`)
	req, _ := http.NewRequest("POST", "/tag", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "test tag" {
		t.Errorf("Expected tag name to be 'test tag'. Got '%v'", m["name"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected tag ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetTag(t *testing.T) {
	clearTable()
	addTags(1)

	req, _ := http.NewRequest("GET", "/tag/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

// main_test.go

func addTags(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO tag(name) VALUES($1)", "Tag "+strconv.Itoa(i))
	}
}

func TestUpdateTag(t *testing.T) {

	clearTable()
	addTags(1)

	req, _ := http.NewRequest("GET", "/tag/1", nil)
	response := executeRequest(req)
	var originalTag map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalTag)

	var jsonStr = []byte(`{"name":"test tag - updated name"}`)
	req, _ = http.NewRequest("PUT", "/tag/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalTag["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalTag["id"], m["id"])
	}

	if m["name"] == originalTag["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalTag["name"], m["name"], m["name"])
	}

}

func TestDeleteTag(t *testing.T) {
	clearTable()
	addTags(1)

	req, _ := http.NewRequest("GET", "/tag/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/tag/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/tag/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestAssignTagToProduct(t *testing.T) {

	clearTable()

	addProducts(1)
	addTags(1)

	var reqStr = "/product/1/tag/1"
	var jsonStr = []byte(`{}`)

	req, _ := http.NewRequest("POST", reqStr, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != 1.0 {
		t.Errorf("Expected tag assignment ID to be '1'. Got '%v'", m["id"])
	}

	if m["productID"] != 1.0 {
		t.Errorf("Expected product ID to be '1'. Got '%v'", m["productID"])
	}

	if m["tagID"] != 1.0 {
		t.Errorf("Expected assigned tag ID to be '1'. Got '%v'", m["tagID"])
	}

}

func addTagAssignment(productID, tagID int) {

	a.DB.Exec("INSERT INTO productToTagAssignment(productID,tagID) VALUES($1,$2)", productID, tagID)
}

func TestDeleteAssignedTagFromProduct(t *testing.T) {

	clearTable()

	addProducts(1)
	addTags(1)
	addTagAssignment(1, 1)

	var reqStr = "/product/1/tag/1"
	var jsonStr = []byte(`{}`)

	req, _ := http.NewRequest("DELETE", reqStr, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

}

func TestCascadeAssignedTagFromProduct(t *testing.T) {

	clearTable()

	addProducts(1)
	addTags(1)
	addTagAssignment(1, 1)

	var reqStr = "/product/1"

	req, _ := http.NewRequest("DELETE", reqStr, nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	reqStr = "/product/1/tag/1"

	req, _ = http.NewRequest("GET", reqStr, nil)

	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)

}

func TestGetAllProductsWithTag(t *testing.T) {

	clearTable()

	addProducts(2)
	addTags(1)
	addTagAssignment(1, 1)
	addTagAssignment(2, 1)

	var reqStr = "/tag/1/products"

	req, _ := http.NewRequest("GET", reqStr, nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

}

func TestGetAllTagsAssignedToProduct(t *testing.T) {

	clearTable()

	addProducts(1)
	addTags(2)
	addTagAssignment(1, 1)
	addTagAssignment(1, 2)

	var reqStr = "/product/1/tags"

	req, _ := http.NewRequest("GET", reqStr, nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

}

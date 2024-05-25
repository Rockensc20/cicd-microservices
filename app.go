// app.go

package main

import (
	"database/sql"
	"fmt"
	"log"

	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(user,port, password, dbname string) {
	//fmt.Print(fmt.Sprintf("user=%s password=%s port=5416 dbname=%s?sslmode=disable", user, password, dbname))
	connectionString :=
		fmt.Sprintf("user=%s password=%s port=%s dbname=%s sslmode=disable", user,port, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()

	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8888", a.Router))
}

func (a *App) getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	p := product{ID: id}
	if err := p.getProduct(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Product not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) getProducts(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	products, err := getProducts(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, products)
}

func (a *App) createProduct(w http.ResponseWriter, r *http.Request) {
	var p product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := p.createProduct(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, p)
}

func (a *App) updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var p product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	p.ID = id

	if err := p.updateProduct(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func (a *App) deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}

	p := product{ID: id}
	if err := p.deleteProduct(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

/*
##################
Tags Functionality
##################
*/

func (a *App) getTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid tag ID")
		return
	}

	t := tag{ID: id}
	if err := t.getTag(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Tag not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, t)
}

func (a *App) getTags(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	products, err := getTags(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, products)
}

func (a *App) createTag(w http.ResponseWriter, r *http.Request) {
	var t tag
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := t.createTag(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, t)
}

func (a *App) updateTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid tag ID")
		return
	}

	var t tag
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	t.ID = id

	if err := t.updateTag(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, t)
}

func (a *App) deleteTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Tag ID")
		return
	}

	t := tag{ID: id}
	if err := t.deleteTag(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// ###############################################################################

func (a *App) getProductToTagAssignment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, errProduct := strconv.Atoi(vars["productID"])
	if errProduct != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Product ID for retriving the specific tag assignment")
		return
	}

	tagID, errTag := strconv.Atoi(vars["tagID"])
	if errTag != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Tag ID for retriving the specific tag assignment")
		return
	}

	pta := productToTagAssignment{ProductID: productID, TagID: tagID}

	if err := pta.getProductToTagAssignment(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Tag assignment to product not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, pta)

}

func (a *App) getProductsWithTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tagID, errProduct := strconv.Atoi(vars["id"])
	if errProduct != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid tag ID for retrieving the products that have it assigned to it")
		return
	}
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	products, err := getProductsWithTagAssigned(a.DB, tagID, start, count)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "No products found with the tag")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, products)
}

func (a *App) getTagsOfProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	productID, errProduct := strconv.Atoi(vars["productID"])
	if errProduct != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID for retrieving the tags assigned to it")
		return
	}

	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	products, err := getTagsAssignedToProduct(a.DB, productID, start, count)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "No tags found on the product")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, products)
}

func (a *App) createProductToTagAssignment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, errProduct := strconv.Atoi(vars["productID"])
	if errProduct != nil {
		respondWithError(w, http.StatusAccepted, "Invalid Product ID")
		return
	}

	tagID, errTag := strconv.Atoi(vars["tagID"])
	if errTag != nil {
		respondWithError(w, http.StatusForbidden, "Invalid Tag ID")
		return
	}

	pta := productToTagAssignment{ProductID: productID, TagID: tagID}

	if err := pta.createProductToTagAssignment(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, pta)
}

func (a *App) deleteProductToTagAssignment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productid, errProduct := strconv.Atoi(vars["productID"])
	tagID, errTag := strconv.Atoi(vars["tagID"])
	if errProduct != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID for deleting the assinged tag")
		return
	}

	if errTag != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid tag ID for deleting the assinged tag")
		return
	}

	pta := productToTagAssignment{ProductID: productid, TagID: tagID}
	if err := pta.deleteProductToTagAssignmentByProductAndTag(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

/*
######
Routes
######
*/

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/product/{productID:[0-9]+}/tags", a.getTagsOfProduct).Methods("GET")
	a.Router.HandleFunc("/product/{productID:[0-9]+}/tag/{tagID:[0-9]+}", a.getProductToTagAssignment).Methods("GET")
	a.Router.HandleFunc("/product/{productID:[0-9]+}/tag/{tagID:[0-9]+}", a.createProductToTagAssignment).Methods("POST")
	a.Router.HandleFunc("/product/{productID:[0-9]+}/tag/{tagID:[0-9]+}", a.deleteProductToTagAssignment).Methods("DELETE")

	a.Router.HandleFunc("/products", a.getProducts).Methods("GET")
	a.Router.HandleFunc("/product", a.createProduct).Methods("POST")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.getProduct).Methods("GET")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.updateProduct).Methods("PUT")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.deleteProduct).Methods("DELETE")

	a.Router.HandleFunc("/tags", a.getTags).Methods("GET")
	a.Router.HandleFunc("/tag", a.createTag).Methods("POST")
	a.Router.HandleFunc("/tag/{id:[0-9]+}", a.getTag).Methods("GET")
	a.Router.HandleFunc("/tag/{id:[0-9]+}/products", a.getProductsWithTag).Methods("GET")
	a.Router.HandleFunc("/tag/{id:[0-9]+}", a.updateTag).Methods("PUT")
	a.Router.HandleFunc("/tag/{id:[0-9]+}", a.deleteTag).Methods("DELETE")

}

package APIServer

import (
	"encoding/json"
	"fmt"
	dbm "main/db"
	"main/models"
	"net/http"
	s "strings"
)

func log(t string, u string) {
	fmt.Println(fmt.Sprintf("[%s]: %s", t, u))
}
func AddItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	log(r.Method, r.URL.String())
	if r.Method != http.MethodPost {
		http.Error(w, "This path takes a POST request only.", http.StatusBadRequest)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Access form values
	product_name := r.Form.Get("product_name")
	if s.TrimSpace(product_name) == "" {
		http.Error(w, "Must have a product name.", http.StatusBadRequest)
		return
	}
	// Do something with the form data
	id, err := dbm.InsertIntoTable(models.Product{Name: product_name})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	dt, err := json.Marshal(models.Product{ID: id, Name: product_name})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(dt)
}
func GetItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method != http.MethodGet {
		http.Error(w, "This path takes a GET request only.", http.StatusBadRequest)
		return
	}
	data, err := dbm.SelectAllFrom(models.Product{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dt, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(dt)
}
func DeleteItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method != http.MethodPost {
		http.Error(w, "This path only takes a POST request.", http.StatusBadRequest)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	id := r.Form.Get("ID")
	_, err = dbm.DeleteWhere(models.Product{ID: id})
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting the product. (%s)", err.Error()), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Success!")
}

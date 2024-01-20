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
func AddToListPOST(w http.ResponseWriter, r *http.Request) {
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
	_, err = dbm.InsertIntoTable(models.Product{Name: product_name})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	fmt.Fprintf(w, "Success!")
}
func GetListGET(w http.ResponseWriter, r *http.Request) {
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

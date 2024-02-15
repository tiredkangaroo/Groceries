package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"main/APIServer"
	dbm "main/db"
	"main/models"
	"net/http"
)

const (
	USER     = "postgres"
	PASSWORD = ""
	DB_NAME  = "groceries_app"
)

func SetupTables() {
	_, err := dbm.CreateTableIfNotExists(models.Product{})
	if err != nil {
		panic(err)
	}
}
func main() {
	db, err := dbm.ConnectToDB(USER, PASSWORD, DB_NAME)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(fmt.Sprintf("Unable to ping the database. Error: %s", err))
	}
	SetupTables()
	http.HandleFunc("/add", APIServer.AddItem)
	http.HandleFunc("/get", APIServer.GetItems)
	http.HandleFunc("/delete", APIServer.DeleteItem)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	fmt.Println("Running API Server at http://[::1]:8080")
	http.ListenAndServe(":8080", nil)
}

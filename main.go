package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/tiredkangaroo/sculpt"
)

type Item struct {
	ID   string
	Name string
}

func main() {
	godotenv.Load(".env")

	user := os.Getenv("POSTGRES_CONNECTION_USER")
	password := os.Getenv("POSTGRES_CONNECTION_PASSWORD")
	dbname := os.Getenv("POSTGRES_CONNECTION_DBNAME")
	err := sculpt.Connect(
		user,
		password,
		dbname,
	)

	if err != nil {
		sculpt.LogError(err.Error())
		return
	}

	item := sculpt.Register(new(Item))
	err = item.Save()
	if err != nil {
		sculpt.LogError(err.Error())
		return
	}

	startAPIServer(item)
}

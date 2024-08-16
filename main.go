package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/tiredkangaroo/sculpt"
)

type Item struct {
	ID          string `json:"id"`
	DateCreated string `json:"date_created"`
	Name        string `json:"name"`
}

func main() {
	godotenv.Load(".env")

	mode := os.Getenv("MODE")
	connectionURI := os.Getenv("POSTGRES_CONNECTION_URI")

	if mode == "" {
		mode = "development"
	}

	err := sculpt.Connect(connectionURI)

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

	startServer(mode, item)
}

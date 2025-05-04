package main

import (
	"log"

	"ualabackend/api"
	"ualabackend/db"

	_ "ualabackend/docs"
)

func main() {
	_, err := db.InitDB()
	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}

	api.InitAPI()
}

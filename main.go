package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/tanvoid0/dev-bot/data"
	"github.com/tanvoid0/dev-bot/server"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	_, err = data.SetupDatabase()
	if err != nil {
		fmt.Println("Error connecting to database")
	}
	server.Run()

}

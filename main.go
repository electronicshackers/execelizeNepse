package main

import (
	"fmt"
	"nepse-backend/api"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("err", err)
	}
	api.Run()
}

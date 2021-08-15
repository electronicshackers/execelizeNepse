package main

import (
	"fmt"
	"nepse-backend/api"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("err", err)
		os.Exit(0)
	}
	api.Run()
}

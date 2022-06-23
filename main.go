package main

import (
	"fmt"
	"log"
	"nepse-backend/api"
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
)

func main() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://74756990a28846ad925c4407371b0a14@o1184422.ingest.sentry.io/6301992",
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	err = godotenv.Load()

	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	if err != nil {
		fmt.Println("err", err)
		os.Exit(0)
	}
	api.Run()
}

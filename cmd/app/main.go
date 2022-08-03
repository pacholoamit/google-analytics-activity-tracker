package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println(os.Getenv("INFO_FOOTTRAFFIK_OAUTH2_CLIENT_AUTH_PROVIDER"))
}

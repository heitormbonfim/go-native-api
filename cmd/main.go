package main

import (
	"log"
	"net/http"

	"github.com/heitormbonfim/go-native-api/config"
)

func main() {
	dbConnection := config.SetupDB()
	defer dbConnection.Close()

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
		return
	}
}

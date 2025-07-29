package main

import (
	"log"
	"net/http"

	"github.com/heitormbonfim/go-native-api/config"
	"github.com/heitormbonfim/go-native-api/models"
)

func main() {
	dbConnection := config.SetupDB()
	defer dbConnection.Close()

	_, err := dbConnection.Exec(models.CreateTableSQL)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
		return
	}
}

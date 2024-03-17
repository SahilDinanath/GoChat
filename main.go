package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/SahilDinanath/GoChat/internal/database"
	"github.com/SahilDinanath/GoChat/internal/routes"
)

func main() {
	fmt.Println("Connecting to database...")
	database.InitDatabaseConnection()
	routes.InitRoutes()
	fmt.Println("running server...")

	log.Fatal(http.ListenAndServe(":8000", nil))
}

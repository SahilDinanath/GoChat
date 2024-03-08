package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	fmt.Println("running server...")

	h1 := func(w http.ResponseWriter, r *http.Request) {
		tml := template.Must(template.ParseFiles("index.html"))
		tml.Execute(w, nil)
	}

	http.HandleFunc("/", h1)

	log.Fatal(http.ListenAndServe(":8000", nil))
}

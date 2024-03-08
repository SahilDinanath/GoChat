package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Message struct {
	Text string
}

func main() {
	fmt.Println("running server...")

	h1 := func(w http.ResponseWriter, r *http.Request) {
		tml := template.Must(template.ParseFiles("index.html"))
		tml.Execute(w, nil)
	}

	sendMessage := func(w http.ResponseWriter, r *http.Request) {
		message := r.PostFormValue("text-bar")
		tml := template.Must(template.ParseFiles("index.html"))
		tml.ExecuteTemplate(w, "message-element", Message{Text: message})
	}

	http.HandleFunc("/", h1)
	http.HandleFunc("/send-message/", sendMessage)

	log.Fatal(http.ListenAndServe(":8000", nil))
}

package routes

import (
	"html/template"
	"log"
	"net/http"

	"github.com/SahilDinanath/GoChat/internal/database"
)

func InitRoutes() {
	http.HandleFunc("/", chatroom)
	http.HandleFunc("/home", homePage)
	http.HandleFunc("/register", registerPage)
	http.HandleFunc("/send-message/", sendMessage)

}

func homePage(w http.ResponseWriter, r *http.Request) {
	tml := template.Must(template.ParseFiles("index.html"))
	tml.Execute(w, nil)
}
func registerPage(w http.ResponseWriter, r *http.Request) {
	tml := template.Must(template.ParseFiles("register.html"))
	tml.Execute(w, nil)
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	message := r.PostFormValue("text")
	/*
		I need to implement saving message to database from here and also call a function to broadcast to all existing chats
	*/
	database.SaveMessage(message, 1)

	data := map[string]string{
		"Text": message,
	}
	tml := template.Must(template.ParseFiles("index.html"))
	tml.ExecuteTemplate(w, "message-element", data)
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	database.SaveUser(r.FormValue("username"), r.FormValue("email"), r.FormValue("password"))

}

func chatroom(w http.ResponseWriter, r *http.Request) {
	tml := template.Must(template.ParseFiles("index.html"))
	messages, err := database.GetMessages(1)

	if err != nil {
		log.Fatal(err)
	}

	data := map[string]any{
		"Messages": messages,
	}

	tml.Execute(w, data)
}

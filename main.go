package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/go-sql-driver/mysql"
)

type Message struct {
	MessageId int
	Username  string
	Text      string
	Timestamp string
}

type Member struct {
	MemberId int
	UserId   int
	RoomId   int
}

func getMessagesFromChat(roomId int) ([]Message, error) {
	var messages []Message

	messageRows, err := db.Query("SELECT message_id, username, message, timestamp FROM messages INNER JOIN members on members.member_id = messages.member_id INNER JOIN users on users.user_id = members.user_id WHERE members.room_id = ?;", roomId)

	fmt.Println("here")
	defer messageRows.Close()

	if err != nil {
		return nil, fmt.Errorf("getMessagesByChat %q: %v", roomId, err)
	}

	for messageRows.Next() {
		var message Message

		if err := messageRows.Scan(&message.MessageId, &message.Username, &message.Text, &message.Timestamp); err != nil {
			return nil, fmt.Errorf("getMessagesByChat %q: %v", roomId, err)
		}
		messages = append(messages, message)
	}

	return messages, nil
}

var db *sql.DB

func main() {
	fmt.Println("running server...")

	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "gochat",
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()

	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connecting to database...")

	h1 := func(w http.ResponseWriter, r *http.Request) {
		tml := template.Must(template.ParseFiles("index.html"))
		tml.Execute(w, nil)
	}

	sendMessage := func(w http.ResponseWriter, r *http.Request) {
		message := r.PostFormValue("text-bar")
		tml := template.Must(template.ParseFiles("index.html"))
		tml.ExecuteTemplate(w, "message-element", Message{Text: message})
	}

	messages, err := getMessagesFromChat(1)

	for _, message := range messages {
		fmt.Println(message.Text)
	}
	http.HandleFunc("/", h1)
	http.HandleFunc("/send-message/", sendMessage)

	log.Fatal(http.ListenAndServe(":8000", nil))
}

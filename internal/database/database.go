package database

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
)

var db *sql.DB

func InitDatabaseConnection() {
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
}

type Message struct {
	MessageId int
	Username  string
	Text      string
	Timestamp string
}

func GetMessages(roomId int) ([]Message, error) {
	var messages []Message

	messageRows, err := db.Query("SELECT message_id, username, message, timestamp FROM messages INNER JOIN members on members.member_id = messages.member_id INNER JOIN users on users.user_id = members.user_id WHERE members.room_id = ?;", roomId)

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

func SaveMessage(message string, memberId int) error {

	_, err := db.Query("INSERT INTO messages(member_id, message) VALUES (?,?)", memberId, message)

	if err != nil {
		return fmt.Errorf("saveMessagesFromChat %q, %q: %v", memberId, message, err)
	}

	return nil
}

func SaveUser(userName string, email string, password string) error {
	_, err := db.Query("INSERT INTO users(username, email, password) VALUES (?,?,?)", userName, email, password)

	if err != nil {
		return fmt.Errorf("SaveUser %q, %q, %q: %v", userName, email, password, err)
	}
	return nil
}

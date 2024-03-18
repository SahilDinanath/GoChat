package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/go-sql-driver/mysql"
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

type User struct {
	UserId   int
	Username string
	Email    string
	Password string
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

func SaveMessage(message string, userId int, roomId int) error {
	var memberId int

	var err error

	err = db.QueryRow("SELECT member_id FROM members WHERE user_id = ? AND room_id = ?", userId, roomId).Scan(&memberId)

	if err != nil {
		return fmt.Errorf("saveMessagesFromChat %q, %q: %v", memberId, message, err)
	}

	_, err = db.Query("INSERT INTO messages(member_id, message) VALUES (?,?)", memberId, message)

	if err != nil {
		return fmt.Errorf("saveMessagesFromChat %q, %q: %v", memberId, message, err)
	}

	return nil
}

func isEmpty(strings ...string) bool {
	for _, str := range strings {
		if str == "" {
			return false
		}
	}
	return true
}
func isValidEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	match, _ := regexp.MatchString(regex, email)

	return match
}

func LoginUser(email string, password string) (*User, error) {
	if !isEmpty(email, password) {
		log.Print(email, password)
		return nil, fmt.Errorf("Fields left empty!")
	}

	if !isValidEmail(email) {
		return nil, fmt.Errorf("Invalid email!")
	}

	var user User

	err := db.QueryRow("Select * from users where email=?", email).Scan(&user.UserId, &user.Username, &user.Email, &user.Password)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("Email not found! User should register.")
	}

	if password != user.Password {
		return nil, fmt.Errorf("Password is incorrect!")
	}

	return &user, nil

}

func SaveUser(userName string, email string, password string) error {

	if !isEmpty(userName, email, password) {
		return fmt.Errorf("Fields left empty!")
	}

	if !isValidEmail(email) {
		return fmt.Errorf("Invalid email!")
	}

	_, err := db.Query("INSERT INTO users(username, email, password) VALUES (?,?,?)", userName, email, password)

	if err != nil {
		return fmt.Errorf("Email already registered!")
	}

	return nil
}

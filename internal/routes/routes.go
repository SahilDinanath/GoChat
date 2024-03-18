package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/SahilDinanath/GoChat/internal/database"
)

func InitRoutes() {
	http.HandleFunc("/", chatroom)
	http.HandleFunc("/home", homePage)
	http.HandleFunc("/register", registerPage)
	http.HandleFunc("/register-user/", registerUser)
	http.HandleFunc("/login/", loginPage)
	http.HandleFunc("/login-user/", loginUser)
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
func loginPage(w http.ResponseWriter, r *http.Request) {
	tml := template.Must(template.ParseFiles("login.html"))
	tml.Execute(w, nil)
}
func sendMessage(w http.ResponseWriter, r *http.Request) {
	message := r.PostFormValue("text")
	/*
		I need to implement saving message to database from here and also call a function to broadcast to all existing chats
	*/
	cookie, err := r.Cookie("user_id")

	if err != nil {
		log.Printf("user cookie error: %v", err)
	}

	userId, err := strconv.Atoi(cookie.Value)
	fmt.Println(userId)
	database.SaveMessage(message, userId, 1)

	data := map[string]string{
		"Text": message,
	}
	tml := template.Must(template.ParseFiles("index.html"))
	tml.ExecuteTemplate(w, "message-element", data)
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	user, err := database.LoginUser(r.FormValue("email"), r.FormValue("password"))

	var htmlContent string

	if err != nil {
		log.Println(err)

		htmlContent = fmt.Sprintf(`
	<div role="alert" class="alert alert-warning">
		<svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
			<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
		</svg>
		<span>
			Error: %s
		</span>
	</div>
		`, err)

	} else {
		cookie := http.Cookie{
			Name:  "user_id",
			Value: strconv.Itoa(user.UserId),
			Path:  "/",
		}
		http.SetCookie(w, &cookie)
		log.Println("cookie sent.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	tml, _ := template.New("t").Parse(htmlContent)

	tml.Execute(w, nil)
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	err := database.SaveUser(r.FormValue("username"), r.FormValue("email"), r.FormValue("password"))

	htmlContent := `
	<div role="alert" class="alert alert-success">
		<svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
		</svg>
		<span>
			Successfully Registered!
		</span>
	</div>`

	if err != nil {
		log.Println(err)

		htmlContent = fmt.Sprintf(`
	<div role="alert" class="alert alert-warning">
		<svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
			<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
		</svg>
		<span>
			Error: %s
		</span>
	</div>
		`, err)
	} else {

		http.Redirect(w, r, "/", http.StatusSeeOther)

	}

	tml, _ := template.New("t").Parse(htmlContent)

	tml.Execute(w, nil)
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

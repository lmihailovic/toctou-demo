package handlers

import (
	"html/template"
	"net/http"

	"toctou-demo/main/middleware"
)

type User struct {
	ID        int
	FirstName string
	LastName  string
	Balance   float64
	IsAdmin   bool
}

type IndexData struct {
	CurrentUser User
	Users       []User
	Success     string
	Error       string
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := middleware.Store.Get(r, "session")
	userID := session.Values["user_id"].(int)

	var current User
	current.ID = userID
	current.FirstName = session.Values["first_name"].(string)
	current.LastName = session.Values["last_name"].(string)
	current.IsAdmin, _ = session.Values["is_admin"].(bool)

	err := db.QueryRow(`SELECT balance FROM users WHERE id=$1`, userID).Scan(&current.Balance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err := db.Query(`SELECT id, first_name, last_name FROM users WHERE id != $1`, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.FirstName, &u.LastName)
		users = append(users, u)
	}

	data := IndexData{
		CurrentUser: current,
		Users:       users,
		Success:     r.URL.Query().Get("success"),
		Error:       r.URL.Query().Get("error"),
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, data)
}

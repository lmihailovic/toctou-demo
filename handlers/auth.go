package handlers

import (
	"database/sql"
	"html/template"
	"net/http"

	"toctou-demo/main/middleware"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		session, _ := middleware.Store.Get(r, "session")
		if uid := session.Values["user_id"]; uid != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		tmpl := template.Must(template.ParseFiles("templates/login.html"))
		tmpl.Execute(w, nil)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	var id int
	var firstName, lastName string
	var isAdmin bool
	err := db.QueryRow(
		`SELECT id, first_name, last_name, is_admin FROM users WHERE email=$1 AND password=$2`,
		email, password,
	).Scan(&id, &firstName, &lastName, &isAdmin)

	if err == sql.ErrNoRows {
		tmpl := template.Must(template.ParseFiles("templates/login.html"))
		tmpl.Execute(w, map[string]string{"Error": "Wrong credentials."})
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session, _ := middleware.Store.Get(r, "session")
	session.Values["user_id"] = id
	session.Values["first_name"] = firstName
	session.Values["last_name"] = lastName
	session.Values["is_admin"] = isAdmin
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := middleware.Store.Get(r, "session")
	session.Values = map[interface{}]interface{}{}
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

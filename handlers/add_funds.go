package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"toctou-demo/main/middleware"
)

type AddFundsData struct {
	CurrentUser User
	Success     string
	Error       string
}

func AddFundsHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := middleware.Store.Get(r, "session")
	userID := session.Values["user_id"].(int)

	var current User
	current.ID = userID
	current.FirstName = session.Values["first_name"].(string)
	current.LastName = session.Values["last_name"].(string)
	current.IsAdmin, _ = session.Values["is_admin"].(bool)

	db.QueryRow(`SELECT balance FROM users WHERE id=$1`, userID).Scan(&current.Balance)

	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("templates/add_funds.html"))
		tmpl.Execute(w, AddFundsData{CurrentUser: current})
		return
	}

	amountStr := r.FormValue("amount")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 || amount > 100000 {
		tmpl := template.Must(template.ParseFiles("templates/add_funds.html"))
		tmpl.Execute(w, AddFundsData{CurrentUser: current, Error: "Amount not in adequate format"})
		return
	}

	_, err = db.Exec(`UPDATE users SET balance = balance + $1 WHERE id = $2`, amount, userID)
	if err != nil {
		tmpl := template.Must(template.ParseFiles("templates/add_funds.html"))
		tmpl.Execute(w, AddFundsData{CurrentUser: current, Error: "Trabsaction not succesful"})
		return
	}

	db.QueryRow(`SELECT balance FROM users WHERE id=$1`, userID).Scan(&current.Balance)
	tmpl := template.Must(template.ParseFiles("templates/add_funds.html"))
	tmpl.Execute(w, AddFundsData{CurrentUser: current, Success: "Succesful transaction!"})
}

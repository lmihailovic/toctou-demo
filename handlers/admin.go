package handlers

import (
	"html/template"
	"net/http"
	"time"
)

type Transfer struct {
	ID             int
	SenderName     string
	SenderLastName string
	RecipName      string
	RecipLastName  string
	Amount         float64
	Timestamp      time.Time
}

type AdminData struct {
	Transfers []Transfer
}

func AdminHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT t.id,
		       s.first_name, s.last_name	,
		       rc.first_name, rc.last_name,
		       t.amount, t.timestamp
		FROM transfers t
		JOIN users s  ON t.sender_id    = s.id
		JOIN users rc ON t.recipient_id = rc.id
		ORDER BY t.timestamp DESC
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var transfers []Transfer
	for rows.Next() {
		var tr Transfer
		rows.Scan(&tr.ID, &tr.SenderName, &tr.SenderLastName,
			&tr.RecipName, &tr.RecipLastName, &tr.Amount, &tr.Timestamp)
		transfers = append(transfers, tr)
	}

	tmpl := template.Must(template.ParseFiles("templates/admin.html"))
	tmpl.Execute(w, AdminData{Transfers: transfers})
}

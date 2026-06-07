package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"toctou-demo/main/middleware"
)

func TransferHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	session, _ := middleware.Store.Get(r, "session")
	senderID := session.Values["user_id"].(int)

	recipientIDStr := r.FormValue("recipient_id")
	amountStr := r.FormValue("amount")

	recipientID, err := strconv.Atoi(recipientIDStr)
	if err != nil || recipientID == senderID {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var senderBalance float64
	err = db.QueryRow(`SELECT balance FROM users WHERE id = $1`, senderID).Scan(&senderBalance)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if senderBalance < amount {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//namerno smo dodali ovde sleep na thread
	//da bi bilo lakse da se dogodi toctou
	// time.Sleep(100 * time.Millisecond)

	_, err = db.Exec(`UPDATE users SET balance = balance - $1 WHERE id = $2`, amount, senderID)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	_, err = db.Exec(`UPDATE users SET balance = balance + $1 WHERE id = $2`, amount, recipientID)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	_, err = db.Exec(
		`INSERT INTO transfers (sender_id, recipient_id, amount, timestamp) VALUES ($1, $2, $3, $4)`,
		senderID, recipientID, amount, time.Now(),
	)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/?success=Transfer successful"), http.StatusSeeOther)
}

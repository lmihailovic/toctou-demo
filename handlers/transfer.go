package handlers

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"toctou-demo/main/middleware"
)

var mu sync.Mutex

func TransferHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// mutex osigurava da samo jedna go rutina moze da udje u ovaj kod u bilo kom datom trenutku
	mu.Lock()
	defer mu.Unlock()

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

	// implementirana SQL transakcija radi ocuvanja konzistentnosti podataka
	tx, err := db.Begin()
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	_, err = tx.Exec(
		`UPDATE users SET balance = balance - $1 WHERE id = $2`, amount, senderID)
	if err != nil {
		tx.Rollback()
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	_, err = tx.Exec(
		`UPDATE users SET balance = balance + $1 WHERE id = $2`, amount, recipientID)
	if err != nil {
		tx.Rollback()
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	_, err = tx.Exec(
		`INSERT INTO transfers (sender_id, recipient_id, amount, timestamp) VALUES ($1, $2, $3, $4)`,
		senderID, recipientID, amount, time.Now(),
	)
	if err != nil {
		tx.Rollback()
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/?success=Transfer successful", http.StatusSeeOther)
}

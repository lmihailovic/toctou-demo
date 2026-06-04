package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"toctou-demo/main/handlers"
	"toctou-demo/main/middleware"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost port=5433 user=postgres password=admin dbname=toctou sslmode=disable"
	}

	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("Database not available: %v", err)
	}
	log.Println("Connection successful.")

	handlers.SetDB(db)

	mux := http.NewServeMux()

	mux.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("styles"))))

	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/logout", handlers.LogoutHandler)

	mux.Handle("/", middleware.RequireAuth(http.HandlerFunc(handlers.IndexHandler)))
	mux.Handle("/transfer", middleware.RequireAuth(http.HandlerFunc(handlers.TransferHandler)))
	mux.Handle("/add-funds", middleware.RequireAuth(http.HandlerFunc(handlers.AddFundsHandler)))
	mux.Handle("/admin", middleware.RequireAdmin(http.HandlerFunc(handlers.AdminHandler)))

	log.Println("Server started: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

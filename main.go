package main

import (
	"log"
	"net/http"

	"mtg-manager/db"
	"mtg-manager/handlers"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if err := db.InitDB("mtg_manager.db"); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.DestroyDB()

	r := mux.NewRouter()
	r.HandleFunc("/users", handlers.UserHandler)
	r.HandleFunc("/cards", handlers.CardHandler)
	// r.HandleFunc("/products/{key}", ProductHandler)
	// r.HandleFunc("/articles", ArticlesHandler)
	http.Handle("/", r)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("could not open server")
	}
}

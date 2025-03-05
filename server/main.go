package main

import (
	"fmt"
	"log"
	"net/http"

	"mtg-manager/server/db"
	"mtg-manager/server/handlers"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if err := db.InitDB("mtg_manager.db"); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.DestroyDB()

	r := mux.NewRouter()
	r.HandleFunc("/login", handlers.LoginHandler)
	r.HandleFunc("/users", handlers.UserHandler)
	r.HandleFunc("/cards", handlers.CardHandler)
	r.HandleFunc("/library/{uid}", handlers.LibraryHandler)
	// r.HandleFunc("/products/{key}", ProductHandler)
	// r.HandleFunc("/articles", ArticlesHandler)
	http.Handle("/", r)

	fmt.Println("wow")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("could not open server")
	}
}

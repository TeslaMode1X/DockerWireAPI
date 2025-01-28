package main

import (
	"github.com/TeslaMode1X/DockerWireAPI/Frontend/handlers"
	"log"
	"net/http"
)

func main() {
	// Маршруты для различных страниц
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/books", handlers.BooksHandler)
	http.HandleFunc("/cart", handlers.CartHandler)

	// Статика для CSS и Bootstrap
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

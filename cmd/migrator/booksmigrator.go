package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
)

type Book struct {
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
	Stock  int     `json:"stock"`
}

func main() {
	dbPath, err := getDbPath()
	if err != nil {
		log.Fatal(fmt.Errorf("error getting db path: %v", err))
	}

	db, err := sql.Open("postgres", dbPath)
	if err != nil {
		log.Fatal(fmt.Errorf("error connecting to db: %v", err))
	}
	defer db.Close()

	isExists, err := checkMigrationDone(db)
	if err != nil {
		log.Fatal(err)
	}
	if isExists {
		fmt.Println("Migration already done.")
		return
	}

	content, err := os.ReadFile("books.json")
	if err != nil {
		log.Fatal(fmt.Errorf("error reading file books.json: %v", err))
	}

	var books []Book

	err = json.Unmarshal(content, &books)
	if err != nil {
		log.Fatal(fmt.Errorf("error JSON parsing: %v", err))
	}

	stmt, err := db.Prepare(`
    INSERT INTO books (title, author, price, stock, created_at) 
    VALUES ($1, $2, $3, $4, $5);
  `)
	if err != nil {
		log.Fatal(fmt.Errorf("error preparing query: %v", err))
	}
	defer stmt.Close()

	for _, book := range books {
		_, err := stmt.Exec(book.Title, book.Author, book.Price, book.Stock, time.Now())
		if err != nil {
			log.Fatal(fmt.Errorf("error inserting data: %v", err))
		}
	}

	fmt.Println("Book migrations have been done successfully!")
}

func checkMigrationDone(db *sql.DB) (bool, error) {
	stmt, err := db.Prepare("SELECT EXISTS(SELECT 1 FROM books)")
	if err != nil {
		return false, fmt.Errorf("error preparing query: %v", err)
	}
	defer stmt.Close()

	var exists bool
	err = stmt.QueryRow().Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking migrations: %v", err)
	}

	return exists, nil
}

func getDbPath() (string, error) {
	var dbUserName, dbPass, dbHost, dbName string

	flag.StringVar(&dbUserName, "db-user-name", "postgres", "Имя пользователя БД")
	flag.StringVar(&dbPass, "db-pass", "postgres", "Пароль пользователя БД")
	flag.StringVar(&dbHost, "db-host", "db:5432", "Адрес БД")
	flag.StringVar(&dbName, "db-name", "postgres", "Название БД")
	flag.Parse()

	if dbName == "" {
		return "", fmt.Errorf("db-name is required")
	}

	dbPath := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbUserName, dbPass, dbHost, dbName)

	return dbPath, nil
}

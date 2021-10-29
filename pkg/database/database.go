package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	log.Println("Initializing SQLITE DB")
	createDB()
}

func createDB() {
	os.Remove("./threat_analyser.db")

	db, err := sql.Open("sqlite3", "./threat_analyser.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	createCmd := `
	create table ip (ip_address TEXT PRIMARY KEY,
					 uuid TEXT,
					 created_at TEXT,
					 updated_at TEXT,
					 response_code TEXT);
	`
	_, err = db.Exec(createCmd)
	if err != nil {
		log.Printf("%q: %s\n", err, createCmd)
		return
	}
}

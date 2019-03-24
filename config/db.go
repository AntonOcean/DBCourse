package config

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

var DB *sql.DB = Database()

func Database() *sql.DB {
	connStr := "postgres://postgres:12345678@localhost:5432/postgres"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Invalid DB config:", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("DB unreachable:", err)
	}
	return db
}
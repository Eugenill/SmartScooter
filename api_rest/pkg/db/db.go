package db

import (
	"database/sql"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/errors"
	"log"
)

var DB *sql.DB

func OpenDB() *sql.DB {
	db, err := sql.Open("postgres", "host=localhost port=5432 dbname=smartscooter user=postgres password=postgres sslmode=disable")
	errors.PanicError(err)
	log.Printf("Connected to the PostgreSQL, Host: localhost, Port: 5432, dbname: smartscooter, user: postgres, password: postgres")
	return db
}

package db

import (
	"database/sql"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/errors"
)

var DB *sql.DB

func OpenDB(db *sql.DB) {
	db, err := sql.Open("postgres", "host=localhost port=5432 dbname=smartscooter user=postgres password=postgres sslmode=disable")
	errors.PanicError(err)
}

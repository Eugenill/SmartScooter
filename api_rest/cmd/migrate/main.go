package main

//USAGE:  go run ./cmd/migrate/main.go
import (
	"context"
	"database/sql"
	"github.com/Eugenill/SmartScooter/api_rest/migrations"

	"github.com/sqlbunny/sqlbunny/runtime/bunny"
)

func main() {
	db, err := sql.Open("postgres", "host=localhost port=5432 dbname=smartscooter user=postgres password=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	ctx = bunny.ContextWithDB(ctx, db)

	err = migrations.Store.Run(ctx)
	if err != nil {
		panic(err)
	}
}

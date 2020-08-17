package db

import (
	"context"
	"database/sql"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/sqlbunny/sqlbunny/runtime/bunny"
	"log"
)

var DB *sql.DB

func OpenDB() *sql.DB {
	db, err := sql.Open("postgres", "host=localhost port=5432 dbname=smartscooter user=postgres password=postgres sslmode=disable")
	errors.PanicError(err)
	log.Printf("Connected to the PostgreSQL, Host: localhost, Port: 5432, dbname: smartscooter, user: postgres, password: postgres")
	return db
}

func GinToContextWithDB(ctx *gin.Context) context.Context {
	ctx2 := context.Background()
	return bunny.ContextWithDB(ctx2, DB)

}

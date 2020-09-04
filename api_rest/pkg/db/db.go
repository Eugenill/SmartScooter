package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/sqlbunny/sqlbunny/runtime/bunny"
	"log"
)

const (
	dbName = "smartscooter"
	dbUser = "postgres"
	dbPass = "postgres"
	dbHost = "localhost"
	dbPort = "5432"
)

var DB *sql.DB

func OpenDB() *sql.DB {
	dbSource := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", dbHost, dbPort, dbName, dbUser, dbPass) //%s for strings, %d for int
	db, err := sql.Open("postgres", dbSource)
	errors.PanicError(err)
	log.Printf("Connected to the PostgreSQL, Host: localhost, Port: 5432, dbname: smartscooter, user: postgres, password: postgres")
	return db
}

func GinToContextWithDB(ctx *gin.Context) context.Context {
	ctx2 := context.Background()
	return bunny.ContextWithDB(ctx2, DB)

}

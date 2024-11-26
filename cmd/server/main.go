package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"aidanwoods.dev/go-paseto"
	"github.com/BigStinko/dave-game-auth/internal/auth"
	"github.com/BigStinko/dave-game-auth/internal/db"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	defer dbConn.Close()
	queries := db.New(dbConn)
	key := paseto.NewV4SymmetricKey() // in production load from environment
	server := auth.NewServer(queries, key)

}

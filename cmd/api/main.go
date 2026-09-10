package main

import (
	"context"
	"log"
	"net/http"

	"github.com/adhamelsaady/digital-wallet/internal/api"
	"github.com/adhamelsaady/digital-wallet/internal/config"
	"github.com/adhamelsaady/digital-wallet/internal/db"
)

func main() {
	cfg , err := config.Load()
	if err != nil {
		log.Fatal("confing eror : %v" , err)
	}
	ctx := context.Background()
	pool , err := db.New(ctx , cfg.DatabaseURL)
	if err != nil {
		log.Fatal("db eror : %v" , err)
	}
	defer pool.Close()
	log.Println(" database connection established")

	router := api.NewRouter()

	log.Printf("listening on port : %s" , cfg.ServerPort)
	
	if err := http.ListenAndServe(cfg.ServerPort, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
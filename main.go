package main

import (
	"log"
	"net/http"
	"os"

	"github.com/takashi002/Go-nami/internal/server"
	"github.com/takashi002/Go-nami/internal/store"
)

func main() {
	addr := ":8080"
	if v := os.Getenv("PORT"); v != "" {
		addr = ":" + v
	}

	st := store.New()
	srv, err := server.New(st)
	if err != nil {
		log.Fatalf("server init: %v", err)
	}

	log.Printf("Nami listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, srv.Handler()); err != nil {
		log.Fatal(err)
	}
}

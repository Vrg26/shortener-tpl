package main

import (
	"github.com/Vrg26/shortener-tpl/internal/app/shortUrl"
	"github.com/Vrg26/shortener-tpl/internal/app/shortUrl/db"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	st := db.NewMemoryStorage()
	service := shortUrl.NewService(st)
	handler := shortUrl.NewHandler(*service)
	handler.RegisterChi(r)
	runServer(r)
}

func runServer(sm *chi.Mux) {
	log.Fatal(http.ListenAndServe(":8080", sm))
}

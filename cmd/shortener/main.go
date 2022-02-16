package main

import (
	"github.com/Vrg26/shortener-tpl/internal/app/shorturl"
	"github.com/Vrg26/shortener-tpl/internal/app/shorturl/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {
	log.Fatal(runServer())
}

func runServer() error {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	st := db.NewMemoryStorage()
	service := shorturl.NewService(st)
	handler := shorturl.NewHandler(*service)
	handler.Register(r)

	return http.ListenAndServe(":8080", r)
}

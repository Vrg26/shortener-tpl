package main

import (
	"github.com/Vrg26/shortener-tpl/internal/app/shorturl"
	"github.com/Vrg26/shortener-tpl/internal/app/shorturl/db"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseUrl       string `env:"BASE_URL" envDefault:"http://localhost"`
}

func main() {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(runServer(&cfg))
}

func runServer(cfg *Config) error {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	st := db.NewMemoryStorage()
	service := shorturl.NewService(st)
	handler := shorturl.NewHandler(*service, cfg.BaseUrl)
	handler.Register(r)

	return http.ListenAndServe(cfg.ServerAddress, r)
}

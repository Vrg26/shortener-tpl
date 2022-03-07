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
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
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

	var st db.Storage
	var err error

	if cfg.FileStoragePath == "" {
		st = db.NewMemoryStorage()
	} else {
		st, _ = db.NewFileStorage(cfg.FileStoragePath)
		if err != nil {
			return err
		}
	}

	service := shorturl.NewService(st)
	handler := shorturl.NewHandler(*service, cfg.BaseURL)
	handler.Register(r)

	return http.ListenAndServe(cfg.ServerAddress, r)
}

package main

import (
	"flag"
	"github.com/Vrg26/shortener-tpl/internal/app/middlewares"
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
	SecretKey       string `env:"SECRET_KEY" envDefault:"secret key"`
}

func main() {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "server address")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "base url")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "file storage path")

	flag.Parse()

	log.Fatal(runServer(&cfg))
}

func runServer(cfg *Config) error {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middlewares.Auth(cfg.ServerAddress))
	r.Use(middlewares.Gzip)

	var st db.Storage
	if cfg.FileStoragePath == "" {
		st = db.NewMemoryStorage()
	} else {
		st = db.NewFileStorage(cfg.FileStoragePath)
	}

	service := shorturl.NewService(st)
	handler := shorturl.NewHandler(*service, cfg.BaseURL)
	handler.Register(r)

	return http.ListenAndServe(cfg.ServerAddress, r)
}

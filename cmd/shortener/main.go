package main

import (
	"context"
	"database/sql"
	"flag"
	"github.com/Vrg26/shortener-tpl/internal/app/middlewares"
	"github.com/Vrg26/shortener-tpl/internal/app/shorturl"
	"github.com/Vrg26/shortener-tpl/internal/app/shorturl/db"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	SecretKey       string `env:"SECRET_KEY" envDefault:"secret key"`
	DataBaseDSN     string `env:"DATABASE_DSN" envDefault:"postgres://test:test@localhost:5432/shorturl?sslmode=disable"`
}

//"result": "http://localhost:8080/315ed11c58ed"
func main() {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "server address")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "base url")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "file storage path")
	flag.StringVar(&cfg.DataBaseDSN, "d", cfg.DataBaseDSN, "database connection string")

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

	var service *shorturl.Service
	if cfg.DataBaseDSN != "" {
		dbPostgres, err := sql.Open("postgres", cfg.DataBaseDSN)

		if err != nil {
			return err
		}
		defer dbPostgres.Close()

		r.Get("/ping", PingDB(dbPostgres))

		st := db.NewPostgresStorage(dbPostgres)

		if err := st.MigrateUp("file://migrations"); err != nil {
			return err
		}

		service = shorturl.NewService(st)

	} else if cfg.FileStoragePath == "" {
		st := db.NewMemoryStorage()
		service = shorturl.NewService(st)
	} else {
		st := db.NewFileStorage(cfg.FileStoragePath)
		service = shorturl.NewService(st)
	}
	handler := shorturl.NewHandler(*service, cfg.BaseURL)
	handler.Register(r)
	return http.ListenAndServe(cfg.ServerAddress, r)
}

func PingDB(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := db.PingContext(ctx); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	}
}

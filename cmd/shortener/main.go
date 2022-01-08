package main

import (
	"github.com/Vrg26/shortener-tpl/internal/app/shortUrl"
	"github.com/Vrg26/shortener-tpl/internal/app/shortUrl/db"
	"log"
	"net/http"
)

func main() {
	r := http.ServeMux{}
	st := db.NewMemoryStorage()
	service := shortUrl.NewService(st)
	handler := shortUrl.NewHandler(*service)
	handler.Register(&r)
	runServer(&r)
}

func runServer(sm *http.ServeMux) {
	log.Fatal(http.ListenAndServe(":8080", sm))
}

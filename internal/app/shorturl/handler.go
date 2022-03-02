package shorturl

import (
	"encoding/json"
	"fmt"
	"github.com/Vrg26/shortener-tpl/internal/app/handlers"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"net/url"
)

var _ handlers.Handler = &handler{}

type handler struct {
	shortURLService Service
}

func NewHandler(service Service) *handler {
	return &handler{shortURLService: service}
}

func (h *handler) Register(r *chi.Mux) {
	r.Get("/{ID}", h.GetURL)
	r.Post("/", h.AddURL)
	r.Post("/api/shorten", h.AddShorten)
}

func (h *handler) GetURL(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "ID")
	if id == "" {
		http.Error(w, "Empty path", http.StatusBadRequest)
		return
	}
	shortURL, err := h.shortURLService.GetByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, shortURL.OriginURL, http.StatusTemporaryRedirect)
}

func (h *handler) AddShorten(w http.ResponseWriter, r *http.Request) {
	var v Shorten

	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if v.Url == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	_, err := url.ParseRequestURI(v.Url)
	if err != nil {
		http.Error(w, "url is invalid", http.StatusBadRequest)
		return
	}

	newID, err := h.shortURLService.Add(v.Url)

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	res, err := json.Marshal(Result{Result: fmt.Sprintf("http://%s/%s", r.Host, newID)})
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func (h *handler) AddURL(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[1:]
	if id != "" {
		http.NotFound(w, r)
		return
	}
	if r.Body == nil {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	b, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	originURL := string(b)
	if originURL == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	_, err = url.ParseRequestURI(originURL)
	if err != nil {
		http.Error(w, "url is invalid", http.StatusBadRequest)
		return
	}

	newID, err := h.shortURLService.Add(originURL)

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("http://%s/%s", r.Host, newID)))
}

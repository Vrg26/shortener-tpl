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
	baseURL         string
}

func NewHandler(service Service, baseURL string) *handler {
	return &handler{shortURLService: service, baseURL: baseURL}
}

func (h *handler) Register(r *chi.Mux) {

	r.Get("/{ID}", h.GetURL)
	r.Post("/", h.AddTextURL)
	r.Post("/api/shorten", h.AddJSONURL)
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

func (h *handler) AddJSONURL(w http.ResponseWriter, r *http.Request) {
	var rBody RequestURL
	if err := json.NewDecoder(r.Body).Decode(&rBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if rBody.URL == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	if _, err := url.ParseRequestURI(rBody.URL); err != nil {
		http.Error(w, "url is invalid", http.StatusBadRequest)
		return
	}

	newID, err := h.shortURLService.Add(rBody.URL)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(ResponseURL{Result: fmt.Sprintf("%s/%s", h.baseURL, newID)})
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func (h *handler) AddTextURL(w http.ResponseWriter, r *http.Request) {
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

	if _, err = url.ParseRequestURI(originURL); err != nil {
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
	w.Write([]byte(fmt.Sprintf("%s/%s", h.baseURL, newID)))
}

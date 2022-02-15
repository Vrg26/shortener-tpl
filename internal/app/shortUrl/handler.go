package shortUrl

import (
	"fmt"
	"github.com/Vrg26/shortener-tpl/internal/app/handlers"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"net/url"
)

var _ handlers.Handler = &handler{}

type handler struct {
	shortUrlService Service
}

func NewHandler(service Service) *handler {
	return &handler{shortUrlService: service}
}
func (h *handler) Register(r *http.ServeMux) {
}

func (h *handler) RegisterChi(r *chi.Mux) {
	r.Get("/", h.GetUrl)
	r.Post("/", h.AddUrl)
}

func (h *handler) routeRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.AddUrl(w, r)
	case http.MethodGet:
		h.GetUrl(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *handler) GetUrl(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[1:]
	if id == "" {
		http.Error(w, "Empty path", http.StatusBadRequest)

		return
	}
	shortUrl, err := h.shortUrlService.GetById(id)
	if err != nil {
		http.NotFound(w, r)

		return
	}
	http.Redirect(w, r, shortUrl.OriginUrl, http.StatusTemporaryRedirect)
}

func (h *handler) AddUrl(w http.ResponseWriter, r *http.Request) {
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

	originUrl := string(b)
	if originUrl == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	_, err = url.ParseRequestURI(originUrl)
	if err != nil {
		http.Error(w, "url is invalid", http.StatusBadRequest)
		return
	}

	newId, err := h.shortUrlService.Add(originUrl)

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("http://%s/%s", r.Host, newId)))
}

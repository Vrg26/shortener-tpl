package shorturl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Vrg26/shortener-tpl/internal/app/handlers"
	"github.com/Vrg26/shortener-tpl/internal/app/shorturl/db"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
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
	r.Get("/api/user/urls", h.GetURLsByUserID)
	r.Post("/", h.AddTextURL)
	r.Post("/api/shorten", h.AddJSONURL)
	r.Post("/api/shorten/batch", h.AddBatchURL)
}

func (h *handler) GetURLsByUserID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	userId := ctx.Value("user").(uint32)
	urls, err := h.shortURLService.GetURLsByUserID(ctx, userId)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	respUrls := make([]RespShortURL, len(urls))
	if len(urls) == 0 {
		resp, _ := json.Marshal(respUrls)
		w.WriteHeader(http.StatusNoContent)
		w.Write(resp)
		return
	}

	for index, url := range urls {
		respUrls[index] = RespShortURL{
			ShortURL:    fmt.Sprintf("%s/%s", h.baseURL, url.ID),
			OriginalURL: url.OriginURL,
		}
	}
	resp, err := json.Marshal(respUrls)

	w.Write(resp)
}

func (h *handler) GetURL(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "ID")
	if id == "" {
		http.Error(w, "Empty path", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	shortURL, err := h.shortURLService.GetByID(ctx, id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, shortURL.OriginURL, http.StatusTemporaryRedirect)
}

func (h *handler) AddBatchURL(w http.ResponseWriter, r *http.Request) {
	userId, _ := r.Context().Value("user").(uint32)
	var rBody []RequestBatchURL
	if err := json.NewDecoder(r.Body).Decode(&rBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var shortUrls []db.ShortURL
	for _, reqUrl := range rBody {
		if reqUrl.OriginalURL == "" {
			http.Error(w, fmt.Sprintf("empty url in the record with id %s", reqUrl.CorrelationId), http.StatusBadRequest)
			return
		}

		if _, err := url.ParseRequestURI(reqUrl.OriginalURL); err != nil {
			http.Error(w, fmt.Sprintf("invalid url in the record with id %s", reqUrl.CorrelationId), http.StatusBadRequest)
			return
		}
		shortUrls = append(shortUrls, db.ShortURL{OriginURL: reqUrl.OriginalURL, CorrelationId: reqUrl.CorrelationId})
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	resultURLs, err := h.shortURLService.AddBatchURL(ctx, shortUrls, userId)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respUrls := make([]ResponseBatchURL, len(resultURLs))

	for index, url := range resultURLs {
		respUrls[index] = ResponseBatchURL{
			ShortURL:      fmt.Sprintf("%s/%s", h.baseURL, url.ID),
			CorrelationId: url.CorrelationId,
		}
	}
	resp, err := json.Marshal(respUrls)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)

}

func (h *handler) AddJSONURL(w http.ResponseWriter, r *http.Request) {
	userId, _ := r.Context().Value("user").(uint32)
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
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	newID, err := h.shortURLService.Add(ctx, rBody.URL, userId)
	if err != nil {
		var pe *pq.Error
		if errors.As(err, &pe) && pgerrcode.IsIntegrityConstraintViolation(string(pe.Code)) {

			newID, err = h.shortURLService.GetByOriginalURL(ctx, rBody.URL)
			if err != nil {
				log.Println(err)
				http.Error(w, "Server error", http.StatusInternalServerError)
				return
			}

			res, err := json.Marshal(RespResultURL{Result: fmt.Sprintf("%s/%s", h.baseURL, newID)})
			if err != nil {
				http.Error(w, "Server error", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusConflict)
			w.Write(res)
			return
		}
		log.Println(err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(RespResultURL{Result: fmt.Sprintf("%s/%s", h.baseURL, newID)})
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func (h *handler) AddTextURL(w http.ResponseWriter, r *http.Request) {
	userId, _ := r.Context().Value("user").(uint32)
	id := r.URL.Path[1:]
	if id != "" {
		http.NotFound(w, r)
		return
	}
	if r.Body == nil {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	originURL := string(b)
	if originURL == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	if _, err = url.ParseRequestURI(originURL); err != nil {
		http.Error(w, "url is invalid", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	newID, err := h.shortURLService.Add(ctx, originURL, userId)
	if err != nil {

		var pe *pq.Error
		if errors.As(err, &pe) && pgerrcode.IsIntegrityConstraintViolation(string(pe.Code)) {

			newID, err = h.shortURLService.GetByOriginalURL(ctx, originURL)
			if err != nil {
				log.Println(err)
				http.Error(w, "Server error", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(fmt.Sprintf("%s/%s", h.baseURL, newID)))
			return
		}

		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("%s/%s", h.baseURL, newID)))
}

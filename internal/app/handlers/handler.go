package handlers

import "net/http"

type Handler interface {
	Register(r *http.ServeMux)
}

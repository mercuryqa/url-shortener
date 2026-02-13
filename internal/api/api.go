package api

import (
	"url/internal/usecases"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	usecases *usecases.UrlShortener
}

func NewHandlers(usecases *usecases.UrlShortener) *Handlers {
	return &Handlers{
		usecases: usecases,
	}
}

func (h *Handlers) InitRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/geturl", h.getUrl)
	r.Post("/generate", h.generateURL)

	return r
}

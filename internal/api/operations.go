package api

import (
	"errors"
	"log"
	"net/http"
	"url/internal/errs"
)

func (h *Handlers) getUrl(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Query().Get("url")
	if shortURL == "" {
		http.Error(w, "missing url param", http.StatusBadRequest)
		return
	}

	url, err := h.usecases.GetUrl(r.Context(), shortURL)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			http.Error(w, "url not found", http.StatusNotFound)
		default:
			log.Printf("get url error: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

func (h *Handlers) generateURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	originalURL := r.URL.Query().Get("url")
	if originalURL == "" {
		http.Error(w, "missing url param", http.StatusBadRequest)
		log.Printf("empty url param")
		return
	}

	shortURL, err := h.usecases.GenerateAndSave(ctx, originalURL)
	if err != nil {
		log.Printf("error GenerateAndSaveUrl %v\n", err)
		switch {
		case errors.Is(err, errs.ErrBadRequest):
			http.Error(w, "bad request", http.StatusBadRequest)
		case errors.Is(err, errs.ErrNotFound):
			http.Error(w, "not found", http.StatusNotFound)
		case errors.Is(err, errs.ErrDatabaseError):
			http.Error(w, "service unavailable", http.StatusServiceUnavailable)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if _, err := w.Write([]byte(shortURL)); err != nil {
		log.Printf("failed to write response: %v", err)
	}

}

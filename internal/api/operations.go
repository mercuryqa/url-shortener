package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

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

	type request struct {
		URL string `json:"url"`
	}

	var req request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		log.Printf("invalid json body: %v", err)
		return
	}
	defer r.Body.Close()

	originalURL := strings.TrimSpace(req.URL)
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

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	resp := struct {
		ShortURL string `json:"short_url"`
	}{
		ShortURL: shortURL,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

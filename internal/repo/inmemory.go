package repo

import (
	"context"
	"fmt"
	"sync"
	"url/internal/errs"
)

type InMemoryRepo struct {
	shortToOriginal map[string]string
	originalToShort map[string]string
	mu              sync.RWMutex
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		shortToOriginal: make(map[string]string),
		originalToShort: make(map[string]string),
	}
}

func (r *InMemoryRepo) GetOriginalUrlByShort(_ context.Context, shortURL string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	originalURL, exists := r.shortToOriginal[shortURL]
	if !exists {
		return "", fmt.Errorf("not found: %w", errs.ErrNotFound)
	}

	return originalURL, nil
}

func (r *InMemoryRepo) SaveUrl(_ context.Context, originalURL string, shortURL string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.shortToOriginal[shortURL]; exists {
		return fmt.Errorf("short url already exists")
	}

	if _, exists := r.shortToOriginal[shortURL]; exists {
		return fmt.Errorf("short url already exists")
	}

	r.shortToOriginal[shortURL] = originalURL
	r.originalToShort[originalURL] = shortURL

	return nil
}

func (r *InMemoryRepo) GetShortByOriginal(_ context.Context, originalURL string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	shortURL, exists := r.originalToShort[originalURL]
	if !exists {
		return "", errs.ErrNotFound
	}

	return shortURL, nil

}

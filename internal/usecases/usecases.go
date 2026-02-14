package usecases

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/url"
	"time"

	"github.com/pkg/errors"

	"url/internal/errs"
)

type RepoUrlShortener interface {
	// GetOriginalUrlByShort returns original url by short url
	GetOriginalUrlByShort(ctx context.Context, shortURL string) (string, error)

	// SaveUrl save short url and original url
	SaveUrl(ctx context.Context, originalURL string, shortURL string) error

	// GetShortByOriginal check existing short url by original url
	GetShortByOriginal(ctx context.Context, originalURL string) (string, error)
}

type UrlShortener struct {
	repo RepoUrlShortener
}

func NewUrlShortener(repo RepoUrlShortener) *UrlShortener {
	return &UrlShortener{
		repo: repo,
	}
}

func (u *UrlShortener) GenerateUrl() (string, error) {

	shortURL, err := u.generateShortURL()
	if err != nil {
		return "", fmt.Errorf("failed generating short url: %w", err)
	}

	return shortURL, nil

}

func (u *UrlShortener) GenerateAndSave(ctx context.Context, originalURL string) (string, error) {
	if !validateUrl(originalURL) {
		return "", errors.Wrap(errs.ErrBadRequest, "It's not URL")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	short, err := u.repo.GetShortByOriginal(ctx, originalURL)
	if err != nil {
		return "", err
	} else if short != "" {
		log.Printf("short url already exists")
		return short, nil
	}

	shortURL, err := u.generateShortURL()
	if err != nil {
		return "", err
	}

	log.Printf("successfully generated")

	if err = u.repo.SaveUrl(ctx, originalURL, shortURL); err != nil {
		log.Printf("failed to save url: %v\n", err)
		return "", err
	}

	return shortURL, nil

}

func (u *UrlShortener) GetUrl(ctx context.Context, shortURL string) (string, error) {
	originalURL, err := u.repo.GetOriginalUrlByShort(ctx, shortURL)
	if err != nil {
		return "", err
	}
	return originalURL, err
}

func validateUrl(originalURL string) bool {

	_, err := url.ParseRequestURI(originalURL)
	return err == nil
}

func (u *UrlShortener) generateShortURL() (string, error) {

	shortURLLength := 10
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

	shortStr := make([]byte, shortURLLength)

	for i := range shortStr {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}

		shortStr[i] = charset[idx.Int64()]
	}
	return string(shortStr), nil
}

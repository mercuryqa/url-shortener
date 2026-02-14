package usecases_test

import (
	"context"
	"errors"
	"testing"

	"url/mocks"

	"github.com/stretchr/testify/assert"
)

func TestGetOriginalUrlByShort(t *testing.T) {
	ctx := context.Background()

	mockRepo := mocks.NewRepoUrlShortener(t)

	shortURL := "k7jnJRQp8e"
	expectedOriginal := "https://example.com"

	mockRepo.On("GetOriginalUrlByShort", ctx, shortURL).Return(expectedOriginal, nil)

	originalURL, err := mockRepo.GetOriginalUrlByShort(ctx, shortURL)

	assert.NoError(t, err)
	assert.Equal(t, expectedOriginal, originalURL)

	mockRepo.AssertExpectations(t)
}

func TestGetOriginalUrlByShortNotFound(t *testing.T) {
	ctx := context.Background()

	mockRepo := mocks.NewRepoUrlShortener(t)

	shortURL := "k7jnJRQp8e"

	mockRepo.On("GetOriginalUrlByShort", ctx, shortURL).Return("", nil)

	originalURL, err := mockRepo.GetOriginalUrlByShort(ctx, shortURL)

	assert.NoError(t, err)
	assert.Equal(t, "", originalURL)

	mockRepo.AssertExpectations(t)
}

func TestSaveUrl(t *testing.T) {
	ctx := context.Background()
	originalURL := "https://example.com"
	shortURL := "k7jnJRQp8e"

	mockRepo := mocks.NewRepoUrlShortener(t)

	mockRepo.On("SaveUrl", ctx, originalURL, shortURL).Return(nil)

	err := mockRepo.SaveUrl(ctx, originalURL, shortURL)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSaveUrlError(t *testing.T) {
	ctx := context.Background()
	originalURL := "https://example.com"
	shortURL := "k7jnJRQp8e"

	mockRepo := mocks.NewRepoUrlShortener(t)

	expectedErr := errors.New("failed to save urls")
	mockRepo.On("SaveUrl", ctx, originalURL, shortURL).Return(expectedErr)

	err := mockRepo.SaveUrl(ctx, originalURL, shortURL)

	assert.EqualError(t, err, "failed to save urls")
	mockRepo.AssertExpectations(t)
}

func TestGetShortByOriginal(t *testing.T) {
	ctx := context.Background()
	originalURL := "https://example.com"
	shortURL := "k7jnJRQp8e"

	mockRepo := mocks.NewRepoUrlShortener(t)

	mockRepo.On("GetShortByOriginal", ctx, originalURL).Return(shortURL, nil)

	result, err := mockRepo.GetShortByOriginal(ctx, originalURL)

	assert.NoError(t, err)
	assert.Equal(t, shortURL, result)
	mockRepo.AssertExpectations(t)
}

func TestGetShortByOriginalError(t *testing.T) {
	ctx := context.Background()
	originalURL := "https://example.com"

	mockRepo := mocks.NewRepoUrlShortener(t)

	expectedErr := errors.New("not found")
	mockRepo.On("GetShortByOriginal", ctx, originalURL).Return("", expectedErr)

	result, err := mockRepo.GetShortByOriginal(ctx, originalURL)

	assert.EqualError(t, err, "not found")
	assert.Equal(t, "", result)
	mockRepo.AssertExpectations(t)
}

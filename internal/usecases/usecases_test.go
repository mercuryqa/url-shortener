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

	shortURL := "abc123DEF_"
	expectedOriginal := "https://example.com/test"

	mockRepo.On("GetOriginalUrlByShort", ctx, shortURL).Return(expectedOriginal, nil)

	originalURL, err := mockRepo.GetOriginalUrlByShort(ctx, shortURL)

	assert.NoError(t, err)
	assert.Equal(t, expectedOriginal, originalURL)

	mockRepo.AssertExpectations(t)
}

func TestGetOriginalUrlByShort_NotFound(t *testing.T) {
	ctx := context.Background()

	mockRepo := mocks.NewRepoUrlShortener(t)

	shortURL := "notexist"

	mockRepo.On("GetOriginalUrlByShort", ctx, shortURL).Return("", errors.New("not found"))

	originalURL, err := mockRepo.GetOriginalUrlByShort(ctx, shortURL)

	assert.Error(t, err)
	assert.Equal(t, "", originalURL)

	mockRepo.AssertExpectations(t)
}

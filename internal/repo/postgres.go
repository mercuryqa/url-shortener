package repo

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"url/internal/errs"
	"url/internal/repo/migrations"
)

type UrlRepo struct {
	db *sql.DB
}

func NewUrlRepo(db *sql.DB) *UrlRepo {

	if db == nil {
		log.Fatal("db is nil in NewUrlRepo")
	}

	migrationsDir := os.Getenv("MIGRATION_DIR")
	if migrationsDir == "" {
		migrationsDir = "internal/repo/migrations"
	}

	mig := migrations.NewMigrator(db, migrationsDir)
	if err := mig.Up(); err != nil {
		log.Fatalf("failed running migrations: %v", err)
	}
	fmt.Println("✅ Migrations applied successfully")

	return &UrlRepo{
		db: db,
	}

}

func (r *UrlRepo) GetOriginalUrlByShort(ctx context.Context, shortURL string) (string, error) {
	var originalURL string

	err := r.db.QueryRowContext(ctx, "SELECT original_url FROM url_table WHERE short_url = $1", shortURL).Scan(&originalURL)

	if err == sql.ErrNoRows {
		return "", errs.ErrNotFound
	}
	if err != nil {
		return "", fmt.Errorf("failed to get original url for short URL %s: %w", shortURL, err)
	}

	return originalURL, nil
}

func (r *UrlRepo) SaveUrl(ctx context.Context, originalURL string, shortURL string) error {
	_, err := r.db.ExecContext(
		ctx, "INSERT INTO url_table (original_url, short_url) VALUES ($1, $2) ON CONFLICT DO NOTHING", originalURL, shortURL)
	if err != nil {
		return fmt.Errorf("failed to save urls %s,  %s: %w", originalURL, shortURL, err)
	}

	return nil
}

func (r *UrlRepo) GetShortByOriginal(ctx context.Context, originalURL string) (string, error) {

	var shortURL string

	err := r.db.QueryRowContext(ctx, "SELECT short_url FROM url_table WHERE original_url = $1", originalURL).Scan(&shortURL)

	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return shortURL, nil
}

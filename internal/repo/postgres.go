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

func (r *UrlRepo) GetOriginalUrlByShort(_ context.Context, shortURL string) (string, error) {
	row, err := r.db.Query("SELECT original_url FROM url_table WHERE short_url = $1", shortURL)
	if err != nil {
		log.Printf("failed running query: %v", err)
	}
	defer row.Close()

	if !row.Next() {
		return "", errs.ErrNotFound
	}

	var resultStr string
	err = row.Scan(&resultStr)

	return resultStr, err
}

func (r *UrlRepo) SaveUrl(_ context.Context, originalURL string, shortURL string) error {
	_, err := r.db.Exec("INSERT INTO url_table (original_url, short_url) VALUES ($1, $2) ON CONFLICT DO NOTHING", originalURL, shortURL)
	return err
}

func (r *UrlRepo) GetShortByOriginal(ctx context.Context, originalURL string) (string, error) {
	var short string

	err := r.db.QueryRow(`SELECT short_url FROM url_table WHERE original_url = $1`, originalURL).Scan(&short)

	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return short, nil
}

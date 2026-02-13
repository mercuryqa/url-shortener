package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"url/internal/api"
	"url/internal/repo"
	"url/internal/usecases"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

const (
	httpPort = ":8082"
)

type App struct {
	server *http.Server
}

func NewApp() *App {
	return &App{}
}

func (a *App) Start(ctx context.Context) error {

	сtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed loading .env file")
	}

	storageType := os.Getenv("STORAGE_TYPE")

	if storageType == "" {
		storageType = "inmemory"
	}

	var repository usecases.RepoUrlShortener

	switch storageType {
	case "inmemory":
		repository = repo.NewInMemoryRepo()

	case "persistent":

		dsn := os.Getenv("DB_URI")
		if dsn == "" {
			log.Fatal("DB_URI is empty")
		}

		db, err := sql.Open("pgx", dsn)
		if err != nil {
			return fmt.Errorf("failed opening db: %w", err)
		}

		if err = db.PingContext(ctx); err != nil {
			return fmt.Errorf("failed connecting to database: %w", err)
		}

		defer db.Close()

		repository = repo.NewUrlRepo(db)
	}

	service := usecases.NewUrlShortener(repository)

	h := api.NewHandlers(service)

	r := chi.NewRouter()
	r.Mount("/api/v1", h.InitRouter())
	a.server = &http.Server{
		Addr:    httpPort,
		Handler: r,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Println("starting server on port ", httpPort)
		if err = a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-stop
	log.Println("shutting down server...")

	if err = a.server.Shutdown(сtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	log.Println("server stopped")

	return nil

}

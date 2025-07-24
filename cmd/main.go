package main

import (
	"context"
	"database/sql"
	"gottd/internal/database"
	"gottd/internal/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const address = ":8080"

func router(queries *database.Queries) http.Handler {
	r := chi.NewRouter()
	r.Use(
		middleware.Logger,
	)

	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	r.NotFound(handlers.NewNotFoundHandler().ServeHTTP)

	return r
}

func main() {
	// Open the database
	db, err := sql.Open("sqlite3", "file:db/database.sqlite3?cache=shared&mode=rw")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	queries := database.New(db)

	runServer(queries)
}

func runServer(queries *database.Queries) {
	server := &http.Server{
		Addr:    address,
		Handler: router(queries),
	}

	// Channel to listen for termination signals
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig
		log.Println("Shutdown signal received, shutting down gracefully...")

		// Create context with timeout for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Shutdown server gracefully
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Failed to shutdown server: %v", err)
		}
	}()

	log.Printf("Server is running on %s\n", address)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", address, err)
	}

	log.Println("Server stopped")
}

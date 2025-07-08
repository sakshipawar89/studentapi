package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/sakshipawar89/StudentApi/internal/config"
	student "github.com/sakshipawar89/StudentApi/internal/http/Handler"
	"github.com/sakshipawar89/StudentApi/internal/storage/sqlite"
)

func main() {
	log.Println("[DEBUG] Loading configuration...")
	cfg := config.MustLoad()

	// Initialize DB
	db, err := sqlite.New(cfg)
	if err != nil {
		log.Fatalf("[ERROR] Failed to initialize storage: %v", err)
	}
	defer db.Db.Close()

	slog.Info("Storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// Set up router
	router := mux.NewRouter()

	// Routes
	router.HandleFunc("/api/students", student.New(db)).Methods("POST")         // Add new student
	router.HandleFunc("/api/students", student.GetList(db)).Methods("GET")       // Get all students
	router.HandleFunc("/api/students/{id}", student.GetById(db)).Methods("GET") // Get student by ID

	// Create HTTP server
	server := &http.Server{
		Addr:    cfg.HTTPServer.Address,
		Handler: router,
	}

	// Run server
	go func() {
		log.Printf("[INFO] Server starting at %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[ERROR] Server failed: %v", err)
		}
	}()

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	log.Println("[INFO] Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("[ERROR] Server shutdown failed: %v", err)
	}

	log.Println("[INFO] Server exited properly")
}

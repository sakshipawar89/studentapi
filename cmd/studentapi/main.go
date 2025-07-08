package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sakshipawar89/StudentApi/internal/config"
	student "github.com/sakshipawar89/StudentApi/internal/http/Handler"
)

func main() {
	log.Println("[DEBUG] Loading configuration...")
	cfg := config.MustLoad()
	log.Printf("[DEBUG] Configuration loaded: %+v\n", cfg)

	// Setup router
	log.Println("[DEBUG] Initializing HTTP router...")
	router := http.NewServeMux()
	router.HandleFunc("/api/students", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[DEBUG] Incoming request: %s %s", r.Method, r.URL.Path)
		if r.Method == http.MethodPost {
			log.Println("[DEBUG] Handling POST /api/students")
			student.New()(w, r) // Call your actual handler
		} else {
			log.Println("[DEBUG] Method not allowed")
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// Setup HTTP server using config
	server := &http.Server{
		Addr:    cfg.HTTPServer.Address,
		Handler: router,
	}

	log.Printf("[INFO] Starting server at %s", server.Addr)

	// Listen for shutdown signals
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Run server in a goroutine
	go func() {
		log.Println("[DEBUG] Server is ready to accept connections")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[ERROR] Server failed: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-done
	log.Println("[INFO] Received shutdown signal, initiating graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("[ERROR] Server forced to shutdown: %v", err)
	}

	log.Println("[INFO] Server exited properly")
}

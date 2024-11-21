package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/as-harudeen/students-api/internal/config"
	"github.com/as-harudeen/students-api/internal/http/handlers/student"
	"github.com/as-harudeen/students-api/internal/storage/sqlite"
)

func main() {
	// load config
	cfg := config.Init()
	// database setup

	storage, err := sqlite.New(cfg)
	if err != nil {
		slog.Error("error")
		log.Fatal(err)
	}

	slog.Info("storage initialized")

	// setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	router.HandleFunc("GET /api/students", student.GetList(storage))

	// setup server
	server := http.Server{
		Addr:    cfg.Addrs,
		Handler: router,
	}

	fmt.Println("server started")

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()

		if err != nil {
			log.Fatal("failed to start server")
		}

	}()

	<-done

	slog.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")
}

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ranta0/rest-and-go/api/v1"
	"github.com/ranta0/rest-and-go/app"
	"github.com/ranta0/rest-and-go/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	app, err := app.NewApp(cfg)
	if err != nil {
		panic(err)
	}
	// Stop the Loggers at shutdown
	defer app.Logger.Close()
	defer app.LoggerAPI.Close()
	api.InitAPI(app)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: app.Router,
	}

	go func() {
		app.Logger.Printf("Server is running on :%d\n", cfg.ServerPort)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			app.Logger.Printf("Error: %v\n", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	app.Logger.Println("Shutting down...")

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		app.Logger.Printf("Error during server shutdown: %v\n", err)
	} else {
		app.Logger.Println("Server gracefully stopped")
	}
}

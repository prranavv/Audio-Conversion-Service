package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/prranavv/peritys_submission/internal/server"
)

func main() {
	log.Println("Starting server on port 8080")
	server := server.Run()
	// Wait for interrupt signal to gracefully shut down the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed:%+v", err)
	}
	log.Println("Server gracefully stopped")

}

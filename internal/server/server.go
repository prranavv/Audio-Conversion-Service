package server

import (
	"log"
	"net/http"

	"github.com/prranavv/peritys_submission/internal/handlers"
)

func Run() *http.Server {
	h := handlers.NewHandler()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: routes(h),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()
	return srv
}

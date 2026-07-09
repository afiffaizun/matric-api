package main

import (
	"fmt"
	"log"
	"net/http"

	"github/afiffaizun/matric-api/internal/handler"
	"github/afiffaizun/matric-api/internal/middleware"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handler.HealthCheck)
	mux.HandleFunc("GET /api/matrix", handler.GetMatrix)
	mux.HandleFunc("POST /api/matrix", handler.CreateMatrix)

	wrapped := middleware.Logging(mux)

	addr := ":8080"
	fmt.Printf("Server running on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, wrapped))
}

package server

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type clusterServer struct{}

func NewClusterServer() *clusterServer {
	return &clusterServer{}
}

type myHandler struct{}

func (m *myHandler) ServeHTTP(http.ResponseWriter, *http.Request) {
	// server http request
}

func (s *clusterServer) Start(port string) {
	mux := http.NewServeMux()

	handler := &myHandler{}
	mux.Handle("/command", handler)
	mux.Handle("/query", handler)

	addr := fmt.Sprintf("127.0.0.1:%v", port)

	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	// Move to separate function
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to listen and serve: %v", err)
	}
}

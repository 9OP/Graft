package server

import (
	"encoding/json"
	"fmt"
	"graft/app/usecase/cluster"
	"log"
	"net/http"
	"time"
)

type clusterServer struct {
	repository cluster.UseCase
}

func NewClusterServer(repository cluster.UseCase) *clusterServer {
	return &clusterServer{repository}
}

// Move to clusterPort
func (s *clusterServer) executeCommand() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Command string `json:"command"`
		}

		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			Badrequest(w, err.Error())
			return
		}

		data, err := s.repository.ExecuteCommand(input.Command)
		if err != nil {
			Badrequest(w, err.Error())
			return
		}

		Render(w, http.StatusOK, data)
	})
}

func (s *clusterServer) Start(port string) {
	mux := http.NewServeMux()

	mux.Handle("/command", s.executeCommand())

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

type apiError struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

func Unauthorized(w http.ResponseWriter, description string) {
	Render(
		w,
		http.StatusUnauthorized,
		apiError{
			Code:        401,
			Description: description,
		},
	)
}

func Badrequest(w http.ResponseWriter, description string) {
	Render(
		w,
		http.StatusBadRequest,
		apiError{
			Code:        400,
			Description: description,
		},
	)
}

func NotFound(w http.ResponseWriter, description string) {
	Render(
		w,
		http.StatusNotFound,
		apiError{
			Code:        404,
			Description: description,
		},
	)
}

func ServerError(w http.ResponseWriter, description string) {
	Render(
		w,
		http.StatusInternalServerError,
		apiError{
			Code:        1000,
			Description: description,
		},
	)
}

func TokenInvalid(w http.ResponseWriter) {
	Unauthorized(w, "Token invalid")
}

func InvalidParameters(w http.ResponseWriter) {
	Badrequest(w, "Parameters are invalid")
}

// Render jsonified data
func Render(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

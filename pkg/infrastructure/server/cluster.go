package server

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"graft/pkg/domain"
	"graft/pkg/services/api"

	log "github.com/sirupsen/logrus"
)

type clusterServer struct {
	repository api.UseCase
}

func NewClusterServer(repository api.UseCase) *clusterServer {
	return &clusterServer{repository}
}

func (s clusterServer) Start(port string) {
	mux := http.NewServeMux()

	mux.HandleFunc("/query/", s.query)
	mux.HandleFunc("/command/", s.command)

	addr := fmt.Sprintf("127.0.0.1:%v", port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  2 * time.Second,
	}
	listenAndServe(srv)
}

func (s clusterServer) query(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	weakConsistency := r.URL.Query().Has("weak_consistency")
	queryEntry := string(b)

	data, err := s.repository.ExecuteQuery(queryEntry, weakConsistency)
	if err != nil {
		switch e := err.(type) {
		case *domain.NotLeaderError:
			http.Redirect(w, r, e.Leader.TargetApi(), http.StatusTemporaryRedirect)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	render(w, data)
}

func (s clusterServer) command(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	commandEntry := string(b)

	data, err := s.repository.ExecuteCommand(commandEntry)
	if err != nil {
		switch e := err.(type) {
		case *domain.NotLeaderError:
			http.Redirect(w, r, e.Leader.TargetApi(), http.StatusTemporaryRedirect)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	render(w, data)
}

func render(w http.ResponseWriter, data []byte) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/text")
	w.Write(data)
}

func listenAndServe(srv *http.Server) {
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to listen and serve: %v", err)
	}
}

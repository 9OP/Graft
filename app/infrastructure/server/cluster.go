package server

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"graft/app/domain/entity"
	"graft/app/usecase/cluster"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type clusterServer struct {
	repository cluster.UseCase
}

func NewClusterServer(repository cluster.UseCase) *clusterServer {
	return &clusterServer{repository}
}

func getBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil

}

func (s *clusterServer) executeCommand(w http.ResponseWriter, r *http.Request) {
	command := r.URL.Query().Get("entry")

	data, err := s.repository.ExecuteCommand(command)

	if err != nil {
		switch e := err.(type) {
		case entity.NotLeaderError:
			http.Redirect(w, r, e.Leader.Target(), http.StatusTemporaryRedirect)
		case entity.TimeoutError:
			http.Error(w, err.Error(), http.StatusGatewayTimeout)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	bytes, err := getBytes(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func (s *clusterServer) executeQuery(w http.ResponseWriter, r *http.Request) {
	command := r.URL.Query().Get("entry")

	data, err := s.repository.ExecuteQuery(command)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bytes, err := getBytes(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func (s *clusterServer) index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.executeQuery(w, r)

	case http.MethodPost:
		s.executeCommand(w, r)

	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
func (s *clusterServer) Start(port string) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.index)

	addr := fmt.Sprintf("127.0.0.1:%v", port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  1 * time.Second,
	}
	listenAndServe(srv)
}

func listenAndServe(srv *http.Server) {
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to listen and serve: %v", err)
	}
}

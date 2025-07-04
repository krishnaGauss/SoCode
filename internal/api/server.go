package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/krishnaGauss/SoCode/internal/models"
	"github.com/krishnaGauss/SoCode/internal/storage"
	"github.com/rs/cors"
)

type Server struct {
	storage  *storage.PostgresStorage
	upgrader websocket.Upgrader
}

func NewServer(storage *storage.PostgresStorage) *Server {
	return &Server{
		storage: storage,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins in development
			},
		},
	}
}

func (s *Server) SetupRoutes() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/api/logs", s.queryLogs).Methods("GET")
	r.HandleFunc("/api/logs/search", s.searchLogs).Methods("POST")
	r.HandleFunc("/api/logs/ws", s.handleWebSocket)
	r.HandleFunc("/health", s.healthCheck).Methods("GET")

	// Serve static files
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/dist/")))

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	return c.Handler(r)
}

func (s *Server) queryLogs(w http.ResponseWriter, r *http.Request) {
	query := models.LogQuery{}

	if startTime := r.URL.Query().Get("start_time"); startTime != "" {
		if t, err := time.Parse(time.RFC3339, startTime); err == nil {
			query.StartTime = &t
		}
	}

	if endTime := r.URL.Query().Get("end_time"); endTime != "" {
		if t, err := time.Parse(time.RFC3339, endTime); err == nil {
			query.EndTime = &t
		}
	}

	if limit := r.URL.Query().Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			query.Limit = l
		}
	}

	if offset := r.URL.Query().Get("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil {
			query.Offset = o
		}
	}

	query.Search = r.URL.Query().Get("search")
	query.Source = r.URL.Query()["source"]
	query.Service = r.URL.Query()["service"]
	query.Host = r.URL.Query()["host"]

	for _, level := range r.URL.Query()["level"] {
		query.Level = append(query.Level, models.LogLevel(level))
	}

	logs, err := s.storage.QueryLogs(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"logs":  logs,
		"count": len(logs),
	})
}

func (s *Server) searchLogs(w http.ResponseWriter, r *http.Request) {
	var query models.LogQuery
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logs, err := s.storage.QueryLogs(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"logs":  logs,
		"count": len(logs),
	})
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// Send real-time logs (simplified implementation)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Get recent logs
			now := time.Now()
			fiveSecondsAgo := now.Add(-5 * time.Second)
			query := models.LogQuery{
				StartTime: &fiveSecondsAgo,
				EndTime:   &now,
				Limit:     10,
			}

			logs, err := s.storage.QueryLogs(query)
			if err != nil {
				continue
			}

			if len(logs) > 0 {
				conn.WriteJSON(map[string]interface{}{
					"type": "new_logs",
					"logs": logs,
				})
			}
		}
	}
}

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

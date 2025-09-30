package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"health-caretaker/internal/models"
	"health-caretaker/internal/monitor"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Handler struct contains dependencies for HTTP handlers
type Handler struct {
	monitor          *monitor.Monitor
	metricsCollector interface {
		GetMetrics() string
	}
}

// NewHandler creates a new handler instance
func NewHandler(m *monitor.Monitor, metrics interface{ GetMetrics() string }) *Handler {
	return &Handler{
		monitor:          m,
		metricsCollector: metrics,
	}
}

// HandleIndex serves the main HTML page
func (h *Handler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

// HandleAPIEndpoints handles REST API for endpoints
func (h *Handler) HandleAPIEndpoints(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		endpoints := h.monitor.GetEndpoints()
		json.NewEncoder(w).Encode(endpoints)

	case "POST":
		var endpoint models.Endpoint
		if err := json.NewDecoder(r.Body).Decode(&endpoint); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		h.monitor.AddEndpoint(&endpoint)
		json.NewEncoder(w).Encode(endpoint)

	case "DELETE":
		vars := mux.Vars(r)
		id := vars["id"]
		h.monitor.RemoveEndpoint(id)
		w.WriteHeader(http.StatusOK)
	}
}

// HandleCheckEndpoint manually triggers a check for a specific endpoint
func (h *Handler) HandleCheckEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	endpoint, exists := h.monitor.GetEndpoint(id)
	if !exists {
		http.Error(w, "Endpoint not found", http.StatusNotFound)
		return
	}

	go func() {
		h.monitor.CheckEndpoint(endpoint)
		h.monitor.BroadcastUpdate(endpoint)
	}()

	w.WriteHeader(http.StatusOK)
}

// HandleWebSocket handles WebSocket connections
func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	upgrader := h.monitor.GetUpgrader()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "WebSocket upgrade error", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	h.monitor.AddClient(conn)
	defer h.monitor.RemoveClient(conn)

	// Send initial data
	endpoints := h.monitor.GetEndpoints()
	for _, endpoint := range endpoints {
		message, err := json.Marshal(endpoint)
		if err != nil {
			continue
		}
		conn.WriteMessage(websocket.TextMessage, message)
	}

	// Keep connection alive
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

// HandleMetrics serves Prometheus-style metrics
func (h *Handler) HandleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")

	if h.metricsCollector != nil {
		metrics := h.metricsCollector.GetMetrics()
		w.Write([]byte(metrics))
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("# Metrics not available\n"))
	}
}

// HandleHealthz serves the health check endpoint
func (h *Handler) HandleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Simple health check - if the server is running, it's healthy
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"service":   "health-caretaker",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HandleReadyz serves the readiness check endpoint
func (h *Handler) HandleReadyz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Readiness check - verify that the monitor is initialized and ready
	endpoints := h.monitor.GetEndpoints()

	response := map[string]interface{}{
		"status":               "ready",
		"timestamp":            time.Now().Unix(),
		"service":              "health-caretaker",
		"endpoints_configured": len(endpoints),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

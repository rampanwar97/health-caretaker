package monitor

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"health-monitoring/internal/models"

	"github.com/gorilla/websocket"
)

// Monitor manages endpoint monitoring
type Monitor struct {
	endpoints map[string]*models.Endpoint
	clients   map[*websocket.Conn]bool
	upgrader  websocket.Upgrader
	mutex     sync.RWMutex
	metricsCallback func(*models.Endpoint) // Callback for metrics updates
}

// NewMonitor creates a new monitor instance
func NewMonitor() *Monitor {
	return &Monitor{
		endpoints: make(map[string]*models.Endpoint),
		clients:   make(map[*websocket.Conn]bool),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// SetMetricsCallback sets the callback function for metrics updates
func (m *Monitor) SetMetricsCallback(callback func(*models.Endpoint)) {
	m.metricsCallback = callback
}

// AddEndpoint adds a new endpoint to monitor
func (m *Monitor) AddEndpoint(endpoint *models.Endpoint) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if endpoint.ID == "" {
		endpoint.ID = fmt.Sprintf("endpoint_%d", time.Now().UnixNano())
	}

	if endpoint.Method == "" {
		endpoint.Method = "GET"
	}

	if endpoint.Interval == 0 {
		endpoint.Interval = 30
	}

	if endpoint.Timeout == 0 {
		endpoint.Timeout = 10
	}

	endpoint.Status = "checking"
	m.endpoints[endpoint.ID] = endpoint
}

// RemoveEndpoint removes an endpoint from monitoring
func (m *Monitor) RemoveEndpoint(id string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.endpoints, id)
}

// GetEndpoints returns all monitored endpoints
func (m *Monitor) GetEndpoints() []*models.Endpoint {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	endpoints := make([]*models.Endpoint, 0, len(m.endpoints))
	for _, endpoint := range m.endpoints {
		endpoints = append(endpoints, endpoint)
	}
	return endpoints
}

// GetEndpoint returns a specific endpoint by ID
func (m *Monitor) GetEndpoint(id string) (*models.Endpoint, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	endpoint, exists := m.endpoints[id]
	return endpoint, exists
}

// CheckEndpoint performs a health check on a single endpoint
func (m *Monitor) CheckEndpoint(endpoint *models.Endpoint) {
	start := time.Now()

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: time.Duration(endpoint.Timeout) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Allow self-signed certificates
			},
		},
	}

	// Create request
	req, err := http.NewRequest(endpoint.Method, endpoint.URL, nil)
	if err != nil {
		endpoint.Status = "down"
		endpoint.Error = fmt.Sprintf("Failed to create request: %v", err)
		endpoint.LastCheck = time.Now()
		return
	}

	// Perform request
	resp, err := client.Do(req)
	responseTime := time.Since(start).Milliseconds()
	
	endpoint.LastCheck = time.Now()
	endpoint.ResponseTime = responseTime
	
	if err != nil {
		endpoint.Status = "down"
		endpoint.Error = err.Error()
		endpoint.StatusCode = 0
	} else {
		endpoint.StatusCode = resp.StatusCode
		endpoint.Error = ""
		
		if resp.StatusCode >= 200 && resp.StatusCode < 400 {
			endpoint.Status = "up"
		} else {
			endpoint.Status = "down"
			endpoint.Error = fmt.Sprintf("HTTP %d", resp.StatusCode)
		}
		
		resp.Body.Close()
	}
	
	// Update metrics if callback is set
	if m.metricsCallback != nil {
		m.metricsCallback(endpoint)
	}
}

// StartMonitoring begins monitoring all endpoints
func (m *Monitor) StartMonitoring(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.mutex.RLock()
			endpoints := make([]*models.Endpoint, 0, len(m.endpoints))
			for _, endpoint := range m.endpoints {
				endpoints = append(endpoints, endpoint)
			}
			m.mutex.RUnlock()

			for _, endpoint := range endpoints {
				// Check if it's time to check this endpoint
				if time.Since(endpoint.LastCheck) >= time.Duration(endpoint.Interval)*time.Second {
					go func(ep *models.Endpoint) {
						m.CheckEndpoint(ep)
						m.broadcastUpdate(ep)
					}(endpoint)
				}
			}
		}
	}
}

// broadcastUpdate sends endpoint status to all connected WebSocket clients
func (m *Monitor) broadcastUpdate(endpoint *models.Endpoint) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	message, err := json.Marshal(endpoint)
	if err != nil {
		log.Printf("Error marshaling endpoint update: %v", err)
		return
	}

	for client := range m.clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("Error sending message to client: %v", err)
			client.Close()
			delete(m.clients, client)
		}
	}
}

// GetUpgrader returns the WebSocket upgrader
func (m *Monitor) GetUpgrader() websocket.Upgrader {
	return m.upgrader
}

// AddClient adds a WebSocket client
func (m *Monitor) AddClient(conn *websocket.Conn) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.clients[conn] = true
}

// RemoveClient removes a WebSocket client
func (m *Monitor) RemoveClient(conn *websocket.Conn) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.clients, conn)
}

// BroadcastUpdate sends endpoint status to all connected WebSocket clients (public method)
func (m *Monitor) BroadcastUpdate(endpoint *models.Endpoint) {
	m.broadcastUpdate(endpoint)
}

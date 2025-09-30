package models

import "time"

// Endpoint represents a monitored endpoint
type Endpoint struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	URL          string            `json:"url"`
	Method       string            `json:"method"`
	Interval     int               `json:"interval"` // in seconds
	Timeout      int               `json:"timeout"`  // in seconds
	LastCheck    time.Time         `json:"lastCheck"`
	Status       string            `json:"status"` // "up", "down", "checking"
	StatusCode   int               `json:"statusCode"`
	ResponseTime int64             `json:"responseTime"` // in milliseconds
	Error        string            `json:"error,omitempty"`
	Labels       map[string]string `json:"labels,omitempty"`     // Additional labels for metrics
	ProbeType    string            `json:"probe_type,omitempty"` // e.g., "livez", "readyz", "healthz"
}

// NewEndpoint creates a new endpoint with default values
func NewEndpoint(name, url, method string, interval, timeout int) *Endpoint {
	return &Endpoint{
		Name:     name,
		URL:      url,
		Method:   method,
		Interval: interval,
		Timeout:  timeout,
		Status:   "checking",
	}
}

// IsHealthy returns true if the endpoint is up
func (e *Endpoint) IsHealthy() bool {
	return e.Status == "up"
}

// GetStatusColor returns the CSS class for the status
func (e *Endpoint) GetStatusColor() string {
	switch e.Status {
	case "up":
		return "status-up"
	case "down":
		return "status-down"
	case "checking":
		return "status-checking"
	default:
		return "status-checking"
	}
}

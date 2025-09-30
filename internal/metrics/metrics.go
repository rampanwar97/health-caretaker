package metrics

import (
	"fmt"
	"sync"
	"time"

	"health-caretaker/internal/models"
)

// MetricsCollector collects and serves metrics
type MetricsCollector struct {
	endpoints map[string]*models.Endpoint
	mutex     sync.RWMutex
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		endpoints: make(map[string]*models.Endpoint),
	}
}

// UpdateEndpoint updates the metrics for an endpoint
func (mc *MetricsCollector) UpdateEndpoint(endpoint *models.Endpoint) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	mc.endpoints[endpoint.ID] = endpoint
}

// RemoveEndpoint removes an endpoint from metrics
func (mc *MetricsCollector) RemoveEndpoint(id string) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	delete(mc.endpoints, id)
}

// GetMetrics returns Prometheus-style metrics
func (mc *MetricsCollector) GetMetrics() string {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	var metrics string

	// Add timestamp
	metrics += fmt.Sprintf("# HELP health_monitoring_timestamp Current timestamp\n")
	metrics += fmt.Sprintf("# TYPE health_monitoring_timestamp gauge\n")
	metrics += fmt.Sprintf("health_monitoring_timestamp %d\n", time.Now().Unix())

	// Add endpoint metrics
	for _, endpoint := range mc.endpoints {
		// Probe success (1 = up, 0 = down) - similar to blackbox exporter
		probeSuccess := 0
		if endpoint.Status == "up" {
			probeSuccess = 1
		}

		// Build labels for probe_success metric
		labels := mc.buildLabels(endpoint)

		metrics += fmt.Sprintf("# HELP probe_success Displays whether the probe was successful\n")
		metrics += fmt.Sprintf("# TYPE probe_success gauge\n")
		metrics += fmt.Sprintf("probe_success{%s} %d\n", labels, probeSuccess)

		// Response time
		metrics += fmt.Sprintf("# HELP probe_duration_seconds Returns how long the probe took to complete in seconds\n")
		metrics += fmt.Sprintf("# TYPE probe_duration_seconds gauge\n")
		metrics += fmt.Sprintf("probe_duration_seconds{%s} %.3f\n",
			labels, float64(endpoint.ResponseTime)/1000.0)

		// Status code
		metrics += fmt.Sprintf("# HELP probe_http_status_code Response HTTP status code\n")
		metrics += fmt.Sprintf("# TYPE probe_http_status_code gauge\n")
		metrics += fmt.Sprintf("probe_http_status_code{%s} %d\n",
			labels, endpoint.StatusCode)

		// Last check timestamp
		metrics += fmt.Sprintf("# HELP probe_last_check_timestamp Last check timestamp\n")
		metrics += fmt.Sprintf("# TYPE probe_last_check_timestamp gauge\n")
		metrics += fmt.Sprintf("probe_last_check_timestamp{%s} %d\n",
			labels, endpoint.LastCheck.Unix())

		// Check interval
		metrics += fmt.Sprintf("# HELP probe_interval_seconds Check interval in seconds\n")
		metrics += fmt.Sprintf("# TYPE probe_interval_seconds gauge\n")
		metrics += fmt.Sprintf("probe_interval_seconds{%s} %d\n",
			labels, endpoint.Interval)
	}

	// Summary metrics
	totalEndpoints := len(mc.endpoints)
	upEndpoints := 0
	downEndpoints := 0

	for _, endpoint := range mc.endpoints {
		if endpoint.Status == "up" {
			upEndpoints++
		} else if endpoint.Status == "down" {
			downEndpoints++
		}
	}

	metrics += fmt.Sprintf("# HELP health_monitoring_total_endpoints Total number of monitored endpoints\n")
	metrics += fmt.Sprintf("# TYPE health_monitoring_total_endpoints gauge\n")
	metrics += fmt.Sprintf("health_monitoring_total_endpoints %d\n", totalEndpoints)

	metrics += fmt.Sprintf("# HELP health_monitoring_up_endpoints Number of healthy endpoints\n")
	metrics += fmt.Sprintf("# TYPE health_monitoring_up_endpoints gauge\n")
	metrics += fmt.Sprintf("health_monitoring_up_endpoints %d\n", upEndpoints)

	metrics += fmt.Sprintf("# HELP health_monitoring_down_endpoints Number of unhealthy endpoints\n")
	metrics += fmt.Sprintf("# TYPE health_monitoring_down_endpoints gauge\n")
	metrics += fmt.Sprintf("health_monitoring_down_endpoints %d\n", downEndpoints)

	return metrics
}

// buildLabels builds the label string for metrics
func (mc *MetricsCollector) buildLabels(endpoint *models.Endpoint) string {
	labels := []string{
		fmt.Sprintf("name=\"%s\"", endpoint.Name),
		fmt.Sprintf("url=\"%s\"", endpoint.URL),
	}

	// Add custom labels from endpoint configuration
	if endpoint.Labels != nil && len(endpoint.Labels) > 0 {
		for key, value := range endpoint.Labels {
			// Escape quotes in label values
			escapedValue := mc.escapeLabelValue(value)
			labels = append(labels, fmt.Sprintf("%s=\"%s\"", key, escapedValue))
		}
	}

	// Join all labels
	result := ""
	for i, label := range labels {
		if i > 0 {
			result += ", "
		}
		result += label
	}

	return result
}

// escapeLabelValue escapes special characters in label values
func (mc *MetricsCollector) escapeLabelValue(value string) string {
	// Replace backslashes and quotes with escaped versions
	escaped := ""
	for _, char := range value {
		switch char {
		case '\\':
			escaped += "\\\\"
		case '"':
			escaped += "\\\""
		case '\n':
			escaped += "\\n"
		case '\r':
			escaped += "\\r"
		case '\t':
			escaped += "\\t"
		default:
			escaped += string(char)
		}
	}
	return escaped
}

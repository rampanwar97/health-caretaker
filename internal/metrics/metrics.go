package metrics

import (
	"fmt"
	"sort"
	"strings"
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

	var b strings.Builder

	// Add timestamp
	b.WriteString("# HELP health_monitoring_timestamp Current timestamp\n")
	b.WriteString("# TYPE health_monitoring_timestamp gauge\n")
	b.WriteString(fmt.Sprintf("health_monitoring_timestamp %d\n", time.Now().Unix()))

	// Add endpoint metrics
	for _, endpoint := range mc.endpoints {
		// Probe success (1 = up, 0 = down) - similar to blackbox exporter
		probeSuccess := 0
		if endpoint.Status == "up" {
			probeSuccess = 1
		}

		// Build labels for probe_success metric
		labels := mc.buildLabels(endpoint)

		b.WriteString("# HELP probe_success Displays whether the probe was successful\n")
		b.WriteString("# TYPE probe_success gauge\n")
		b.WriteString(fmt.Sprintf("probe_success{%s} %d\n", labels, probeSuccess))

		// Response time
		b.WriteString("# HELP probe_duration_seconds Returns how long the probe took to complete in seconds\n")
		b.WriteString("# TYPE probe_duration_seconds gauge\n")
		b.WriteString(fmt.Sprintf("probe_duration_seconds{%s} %.3f\n", labels, float64(endpoint.ResponseTime)/1000.0))

		// Status code
		b.WriteString("# HELP probe_http_status_code Response HTTP status code\n")
		b.WriteString("# TYPE probe_http_status_code gauge\n")
		b.WriteString(fmt.Sprintf("probe_http_status_code{%s} %d\n", labels, endpoint.StatusCode))

		// Last check timestamp
		b.WriteString("# HELP probe_last_check_timestamp Last check timestamp\n")
		b.WriteString("# TYPE probe_last_check_timestamp gauge\n")
		b.WriteString(fmt.Sprintf("probe_last_check_timestamp{%s} %d\n", labels, endpoint.LastCheck.Unix()))

		// Check interval
		b.WriteString("# HELP probe_interval_seconds Check interval in seconds\n")
		b.WriteString("# TYPE probe_interval_seconds gauge\n")
		b.WriteString(fmt.Sprintf("probe_interval_seconds{%s} %d\n", labels, endpoint.Interval))
	}

	// Summary metrics
	totalEndpoints := len(mc.endpoints)
	upEndpoints := 0
	downEndpoints := 0

	for _, endpoint := range mc.endpoints {
		switch endpoint.Status {
		case "up":
			upEndpoints++
		case "down":
			downEndpoints++
		}
	}

	b.WriteString("# HELP health_monitoring_total_endpoints Total number of monitored endpoints\n")
	b.WriteString("# TYPE health_monitoring_total_endpoints gauge\n")
	b.WriteString(fmt.Sprintf("health_monitoring_total_endpoints %d\n", totalEndpoints))

	b.WriteString("# HELP health_monitoring_up_endpoints Number of healthy endpoints\n")
	b.WriteString("# TYPE health_monitoring_up_endpoints gauge\n")
	b.WriteString(fmt.Sprintf("health_monitoring_up_endpoints %d\n", upEndpoints))

	b.WriteString("# HELP health_monitoring_down_endpoints Number of unhealthy endpoints\n")
	b.WriteString("# TYPE health_monitoring_down_endpoints gauge\n")
	b.WriteString(fmt.Sprintf("health_monitoring_down_endpoints %d\n", downEndpoints))

	return b.String()
}

// buildLabels builds the label string for metrics
func (mc *MetricsCollector) buildLabels(endpoint *models.Endpoint) string {
	labels := []string{
		fmt.Sprintf("name=\"%s\"", endpoint.Name),
		fmt.Sprintf("url=\"%s\"", endpoint.URL),
	}

	// Add custom labels from endpoint configuration (stable order)
	if len(endpoint.Labels) > 0 {
		keys := make([]string, 0, len(endpoint.Labels))
		for k := range endpoint.Labels {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, key := range keys {
			value := endpoint.Labels[key]
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

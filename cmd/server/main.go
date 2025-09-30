package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"health-caretaker/internal/config"
	"health-caretaker/internal/handlers"
	"health-caretaker/internal/metrics"
	"health-caretaker/internal/models"
	"health-caretaker/internal/monitor"
	"health-caretaker/internal/server"
	"health-caretaker/pkg/logger"
	"health-caretaker/pkg/middleware"

	"github.com/gorilla/mux"
)

// Build-time variables (set via ldflags)
var (
	Version   = "dev"
	CommitSHA = "unknown"
	BuildDate = "unknown"
)

func main() {
	// Parse command line flags
	var (
		configFile  = flag.String("config", "config.json", "Configuration file path")
		showVersion = flag.Bool("version", false, "Show version information")
	)
	flag.Parse()

	// Show version if requested
	if *showVersion {
		println(fmt.Sprintf("Version: %s, BuildDate: %s, CommitSHA: %s, GoVersion: %s",
			Version, BuildDate, CommitSHA, runtime.Version()))
		os.Exit(0)
	}

	// Initialize logger
	log := logger.New()
	log.Info("Starting Health Monitoring Service v%s (commit: %s, built: %s)", Version, CommitSHA, BuildDate)

	// Load configuration
	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		log.Fatal("Failed to load config: %v", err)
	}

	log.Info("Loaded configuration from %s", *configFile)
	log.Info("Found %d endpoints in configuration", len(cfg.Endpoints))

	// Create monitor instance
	monitor := monitor.NewMonitor()

	// Create metrics collector
	metricsCollector := metrics.NewMetricsCollector()

	// Set up metrics callback
	monitor.SetMetricsCallback(func(endpoint *models.Endpoint) {
		metricsCollector.UpdateEndpoint(endpoint)
	})

	// Create handler instance
	handler := handlers.NewHandler(monitor, metricsCollector)

	// Load endpoints from configuration
	for _, endpointConfig := range cfg.Endpoints {
		endpoint := endpointConfig.ToEndpoint()
		monitor.AddEndpoint(endpoint)
		log.Info("Added endpoint: %s (%s)", endpoint.Name, endpoint.URL)
		if endpoint.Labels != nil && len(endpoint.Labels) > 0 {
			log.Info("  Labels: %v", endpoint.Labels)
		} else {
			log.Info("  No custom labels configured")
		}
	}

	// Start monitoring in background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go monitor.StartMonitoring(ctx)

	// Setup main server routes
	mainRouter := mux.NewRouter()

	// Add middleware
	mainRouter.Use(middleware.LoggingMiddleware(log))
	mainRouter.Use(middleware.CORSMiddleware())
	mainRouter.Use(middleware.SecurityMiddleware())

	// Static files
	mainRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// Main page
	mainRouter.HandleFunc("/", handler.HandleIndex)

	// API routes
	api := mainRouter.PathPrefix("/api").Subrouter()
	api.HandleFunc("/endpoints", handler.HandleAPIEndpoints).Methods("GET", "POST")
	api.HandleFunc("/endpoints/{id}", handler.HandleAPIEndpoints).Methods("DELETE")
	api.HandleFunc("/endpoints/{id}/check", handler.HandleCheckEndpoint).Methods("POST")

	// WebSocket
	mainRouter.HandleFunc("/ws", handler.HandleWebSocket)

	// Health check endpoints
	mainRouter.HandleFunc("/healthz", handler.HandleHealthz)
	mainRouter.HandleFunc("/readyz", handler.HandleReadyz)
	log.Info("Health check endpoints enabled at /healthz and /readyz")

	// Create servers
	mainServer := server.New(":"+cfg.Server.Port, mainRouter, "Main", log)

	var metricsServer *server.Server
	if cfg.Metrics.Enabled {
		// Create metrics-only router
		metricsRouter := mux.NewRouter()
		metricsRouter.Use(middleware.LoggingMiddleware(log))
		metricsRouter.HandleFunc(cfg.Metrics.Path, handler.HandleMetrics)
		metricsRouter.HandleFunc("/healthz", handler.HandleHealthz)
		metricsRouter.HandleFunc("/readyz", handler.HandleReadyz)

		metricsServer = server.New(":"+cfg.Metrics.Port, metricsRouter, "Metrics", log)
		log.Info("Metrics available at http://localhost:%s%s", cfg.Metrics.Port, cfg.Metrics.Path)
	}

	// Start metrics server in background if enabled
	if metricsServer != nil {
		if err := metricsServer.StartBackground(); err != nil {
			log.Fatal("Failed to start metrics server: %v", err)
		}
	}

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start main server in background
	go func() {
		if err := mainServer.StartBackground(); err != nil {
			log.Fatal("Failed to start main server: %v", err)
		}
	}()

	log.Info("Health monitoring service started successfully")
	log.Info("Web UI available at http://localhost:%s", cfg.Server.Port)

	// Wait for shutdown signal
	<-sigChan
	log.Info("Shutdown signal received, stopping servers...")

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Stop servers gracefully
	if metricsServer != nil {
		if err := metricsServer.Stop(shutdownCtx); err != nil {
			log.Error("Failed to stop metrics server: %v", err)
		}
	}

	if err := mainServer.Stop(shutdownCtx); err != nil {
		log.Error("Failed to stop main server: %v", err)
	}

	log.Info("Health monitoring service stopped")
}

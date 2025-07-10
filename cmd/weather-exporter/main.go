package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/weather-exporter/internal/config"
	"github.com/weather-exporter/internal/metrics"
	"github.com/weather-exporter/internal/server"
	"github.com/weather-exporter/internal/weather"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "configs/config.yaml", "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	setupLogging(cfg.Logging)

	log.Printf("Starting weather exporter...")

	// Create weather provider
	weatherProvider := weather.NewOpenWeatherMapClient(cfg.Weather.APIKey)

	// Extract city names from config
	cities := make([]string, 0, len(cfg.Cities))
	for _, city := range cfg.Cities {
		cities = append(cities, city.Name)
	}

	// Create metrics collector
	collector := metrics.NewCollector(weatherProvider, cities, cfg.Scraping.Interval)

	// Start collector in background
	ctx, cancel := context.WithCancel(context.Background())
	go collector.Start(ctx)

	// Create and start HTTP server
	srv := server.NewServer(cfg.Server.Port, weatherProvider)
	
	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down...")
		cancel()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		if err := srv.Stop(shutdownCtx); err != nil {
			log.Printf("Error shutting down server: %v", err)
		}
		os.Exit(0)
	}()

	// Start server
	log.Printf("Server listening on port %d", cfg.Server.Port)
	log.Printf("Metrics available at http://localhost:%d/metrics", cfg.Server.Port)
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupLogging(cfg config.LoggingConfig) {
	// For simplicity, using standard log package
	// In production, you would use zap or logrus with proper configuration
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	
	if cfg.Output == "stdout" {
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(os.Stderr)
	}
}
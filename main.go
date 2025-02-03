package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	endpoint string
	apikey   string
	listen   string
)

func main() {
	flag.StringVar(&endpoint, "endpoint", "", "cPanel endpoint")
	flag.StringVar(&apikey, "apikey", "", "API key")
	flag.StringVar(&listen, "listen", "0.0.0.0:8080", "Address and port to listen on (format: ip:port)")

	flag.Parse()

	// Check if required flags are provided
	if endpoint == "" || apikey == "" {
		fmt.Println("Error: --endpoint and --apikey are required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Initialize Prometheus metrics
	initEmailStats()
	initEmailLogs()

	// Handle the /metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<h1>Mail cPanel Exporter</h1>")
		fmt.Fprintln(w, "<p>Visit <a href=\"/metrics\">/metrics</a> for Prometheus metrics.</p>")
	})

	// Start a goroutine to periodically fetch and update metrics
	go updateMetrics()

	// Start the HTTP server
	fmt.Printf("Starting server on %s\n", listen)
	if err := http.ListenAndServe(listen, nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

// Periodically update metrics every 5 minutes
func updateMetrics() {
	log.Println("🔄 Initial fetch of email metrics...")
	fetchEmailTrackUserStats()
	fetchEmailTrackSearch()

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("🔄 Refreshing email metrics...")
		fetchEmailTrackUserStats()
		fetchEmailTrackSearch()
	}
}

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"cpanel_mail_exporter/utils"
)

var (
	endpoint      string
	apikey        string
	listen        string
	webConfigFile string
)

func main() {
	flag.StringVar(&endpoint, "endpoint", "", "cPanel endpoint")
	flag.StringVar(&apikey, "apikey", "", "API key")
	flag.StringVar(&listen, "listen", ":9197", "Address and port to listen on (format: ip:port)")
	flag.StringVar(&webConfigFile, "web.config.file", "", "Path to web configuration YAML file")

	flag.Parse()

	// Check if required flags are provided
	if endpoint == "" || apikey == "" {
		fmt.Println("Error: --endpoint and --apikey are required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Load web configuration if provided
	var webConfig *utils.WebConfig
	var err error
	if webConfigFile != "" {
		webConfig, err = utils.LoadWebConfig(webConfigFile)
		if err != nil {
			log.Fatalf("Error loading web config: %v", err)
		}
	}

	// Initialize Prometheus metrics
	initEmailStats()
	initEmailLogs()

	// Handle the /metrics endpoint with optional Basic Auth
	metricsHandler := promhttp.Handler()
	if webConfig != nil && len(webConfig.BasicAuthUsers) > 0 {
		metricsHandler = utils.BasicAuth(metricsHandler, webConfig.BasicAuthUsers)
	}
	http.Handle("/metrics", metricsHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<h1>Mail cPanel Exporter</h1>")
		fmt.Fprintln(w, "<p>Visit <a href=\"/metrics\">/metrics</a> for Prometheus metrics.</p>")
	})

	// Start a goroutine to periodically fetch and update metrics
	go updateMetrics()

	// Start the HTTP server with optional TLS
	log.Println("Starting server...")
	if webConfig != nil && webConfig.TLSServerConfig.CertFile != "" && webConfig.TLSServerConfig.KeyFile != "" {
		log.Println("TLS is enabled.")
		if err := http.ListenAndServeTLS(listen, webConfig.TLSServerConfig.CertFile, webConfig.TLSServerConfig.KeyFile, nil); err != nil {
			log.Fatalf("Error starting TLS server: %v", err)
		}
	} else {
		log.Println("TLS is disabled, starting without encryption.")
		if err := http.ListenAndServe(listen, nil); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
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

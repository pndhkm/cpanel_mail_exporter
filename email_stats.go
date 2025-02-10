package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"cpanel_mail_exporter/utils"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	emailStats *prometheus.GaugeVec
	emailFlags *prometheus.GaugeVec
)

// Struct for API response
type EmailTrackUserStatsResponse struct {
	Metadata struct {
		Command string `json:"command"`
		Reason  string `json:"reason"`
		Result  int    `json:"result"`
		Version int    `json:"version"`
	} `json:"metadata"`
	Data struct {
		Records []struct {
			DeferCount     int    `json:"DEFERCOUNT"`
			DeferFailCount int    `json:"DEFERFAILCOUNT"`
			Domain         string `json:"DOMAIN"`
			FailCount      int    `json:"FAILCOUNT"`
			Owner          string `json:"OWNER"`
			PrimaryDomain  string `json:"PRIMARY_DOMAIN"`
			SendCount      int    `json:"SENDCOUNT"`
			SuccessCount   int    `json:"SUCCESSCOUNT"`
			TotalSize      int    `json:"TOTALSIZE"`
			User           string `json:"USER"`
		} `json:"records"`
	} `json:"data"`
}

// Initialize Prometheus metrics
func initEmailStats() {
	emailStats = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "email_user_stats",
			Help: "Email statistics per user",
		},
		[]string{"domain", "primary_domain", "user", "owner", "stat_type"},
	)
	emailFlags = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "email_flag_stats",
			Help: "Email flag statistics",
		},
		[]string{"domain", "user", "flag"},
	)
	prometheus.MustRegister(emailStats)
	prometheus.MustRegister(emailFlags)
}

// Fetch email statistics dynamically
func fetchEmailTrackUserStats() {
	startTime, endTime := utils.GetStartAndEndOfDay()
	if startTime == 0 || endTime == 0 {
		log.Println("[❌ Email Stats] Failed to get start and end time.")
		return
	}

	queryParams := url.Values{}
	queryParams.Add("api.version", "1")
	queryParams.Add("starttime", strconv.FormatInt(startTime, 10))
	queryParams.Add("endtime", strconv.FormatInt(endTime, 10))

	apiURL := fmt.Sprintf("https://%s/json-api/emailtrack_user_stats?%s", endpoint, queryParams.Encode())

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Println("[❌ Email Stats] Error creating request:", err)
		return
	}

	req.Header.Add("Authorization", fmt.Sprintf("whm root:%s", apikey))

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("[❌ Email Stats] Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// **Check HTTP response status**
	if resp.StatusCode != http.StatusOK {
		log.Printf("[❌ Email Stats] API request failed! Status: %d - %s\n", resp.StatusCode, http.StatusText(resp.StatusCode))
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("[❌ Email Stats] Error reading response:", err)
		return
	}

	// **Validate JSON response**
	var data EmailTrackUserStatsResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println("[❌ Email Stats] Failed to parse API response. Is the endpoint correct?")
		log.Println("👉 Raw Response:", string(body)) // Debugging hint
		return
	}

	// 🚀 Reset old metrics before updating
	emailStats.Reset()

	for _, record := range data.Data.Records {
		emailStats.WithLabelValues(record.Domain, record.PrimaryDomain, record.User, record.Owner, "defer_count").Set(float64(record.DeferCount))
		emailStats.WithLabelValues(record.Domain, record.PrimaryDomain, record.User, record.Owner, "defer_fail_count").Set(float64(record.DeferFailCount))
		emailStats.WithLabelValues(record.Domain, record.PrimaryDomain, record.User, record.Owner, "fail_count").Set(float64(record.FailCount))
		emailStats.WithLabelValues(record.Domain, record.PrimaryDomain, record.User, record.Owner, "send_count").Set(float64(record.SendCount))
		emailStats.WithLabelValues(record.Domain, record.PrimaryDomain, record.User, record.Owner, "success_count").Set(float64(record.SuccessCount))
		emailStats.WithLabelValues(record.Domain, record.PrimaryDomain, record.User, record.Owner, "total_size").Set(float64(record.TotalSize))
	}

	log.Println("✅ Successfully updated email stats metrics.")
}

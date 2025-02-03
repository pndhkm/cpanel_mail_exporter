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

	"github.com/pndhkm/cpanel_mail_exporter/utils"
	"github.com/prometheus/client_golang/prometheus"
)

var emailLogs *prometheus.GaugeVec

type EmailTrackSearchResponse struct {
	Metadata struct {
		Version    int    `json:"version"`
		Overflowed int    `json:"overflowed"`
		Result     int    `json:"result"`
		Reason     string `json:"reason"`
		Command    string `json:"command"`
	} `json:"metadata"`
	Data struct {
		Records []struct {
			Recipient         string `json:"recipient"`
			Email             string `json:"email"`
			SendUnixTime      int64  `json:"sendunixtime"`
			DeliveryUser      string `json:"deliveryuser"`
			TransportIsRemote int    `json:"transport_is_remote"`
			Host              string `json:"host"`
			Sender            string `json:"sender"`
			Router            string `json:"router"`
			ActionTime        string `json:"actiontime"`
			Message           string `json:"message"`
			DeliveryDomain    string `json:"deliverydomain"`
			MsgID             string `json:"msgid"`
			Size              int    `json:"size"`
			Type              string `json:"type"`
			DeliveredTo       string `json:"deliveredto"`
			Transport         string `json:"transport"`
			User              string `json:"user"`
			SenderHost        string `json:"senderhost"`
			SenderIP          string `json:"senderip"`
			ActionUnixTime    int64  `json:"actionunixtime"`
			SenderAuth        string `json:"senderauth"`
			IP                string `json:"ip"`
			Domain            string `json:"domain"`
		} `json:"records"`
	} `json:"data"`
}

func initEmailLogs() {
	emailLogs = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "email_logs",
			Help: "Email logs from cPanel API",
		},
		[]string{
			"time", "actionunixtime", "delivered_to", "delivery_domain", "delivery_user",
			"domain", "host", "ip", "message", "msg_id", "recipient", "router", "sender",
			"sender_auth", "sender_host", "sender_ip", "sendunixtime", "size", "transport",
			"transport_is_remote", "type", "user",
		},
	)
	prometheus.MustRegister(emailLogs)
}

func fetchEmailTrackSearch() {
	startTime, endTime := utils.GetStartAndEndOfDay()
	if startTime == 0 || endTime == 0 {
		log.Println("[❌ Email Logs] Failed to get start and end time.")
		return
	}

	queryParams := url.Values{}
	queryParams.Add("api.version", "1")
	queryParams.Add("api.filter.enable", "1")
	queryParams.Add("api.filter.a.field", "sendunixtime")
	queryParams.Add("api.filter.a.arg0", strconv.FormatInt(startTime, 10))
	queryParams.Add("api.filter.a.type", "gt")
	queryParams.Add("api.filter.b.field", "sendunixtime")
	queryParams.Add("api.filter.b.arg0", strconv.FormatInt(endTime, 10))
	queryParams.Add("api.filter.b.type", "lt_equal")
	queryParams.Add("api.sort.enable", "1")
	queryParams.Add("api.sort.a.field", "sendunixtime")
	queryParams.Add("api.sort.a.reverse", "0")
	queryParams.Add("success", "1")
	queryParams.Add("failure", "1")
	queryParams.Add("inprogress", "0")
	queryParams.Add("defer", "0")
	queryParams.Add("max_results_by_type", "0")

	apiURL := fmt.Sprintf("https://%s/json-api/emailtrack_search?%s", endpoint, queryParams.Encode())

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Println("[❌ Email Logs] Error creating HTTP request:", err)
		return
	}

	req.Header.Add("Authorization", fmt.Sprintf("whm root:%s", apikey))

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("[❌ Email Logs] Error sending HTTP request:", err)
		return
	}
	defer resp.Body.Close()

	// **Check HTTP response status**
	if resp.StatusCode != http.StatusOK {
		log.Printf("[❌ Email Logs] API request failed! Status: %d - %s\n", resp.StatusCode, http.StatusText(resp.StatusCode))
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("[❌ Email Logs] Error reading API response:", err)
		return
	}

	// **Validate JSON response**
	var data EmailTrackSearchResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println("[❌ Email Logs] Failed to parse API response. Make sure the endpoint is correct.")
		log.Println("👉 Raw Response:", string(body)) // Debugging hint
		return
	}

	// 🚀 Reset old metrics before updating
	emailLogs.Reset()

	for _, record := range data.Data.Records {
		emailLogs.WithLabelValues(
			record.ActionTime,
			strconv.FormatInt(record.ActionUnixTime, 10),
			record.DeliveredTo,
			record.DeliveryDomain,
			record.DeliveryUser,
			record.Domain,
			record.Host,
			record.IP,
			record.Message,
			record.MsgID,
			record.Recipient,
			record.Router,
			record.Sender,
			record.SenderAuth,
			record.SenderHost,
			record.SenderIP,
			strconv.FormatInt(record.SendUnixTime, 10),
			strconv.Itoa(record.Size),
			record.Transport,
			strconv.Itoa(record.TransportIsRemote),
			record.Type,
			record.User,
		).Set(1)
	}

	log.Println("✅ Successfully updated email log metrics.")
}

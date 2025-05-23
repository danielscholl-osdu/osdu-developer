package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Configuration options with environment variable defaults
var (
	port           = getEnvOrDefault("PORT", "9002")
	logLevel       = getEnvOrDefault("LOG_LEVEL", "info")
	defaultLatency = getEnvIntOrDefault("DEFAULT_LATENCY_MS", 0)
)

// ExtAuthRequest represents the ext_authz request structure
type ExtAuthRequest struct {
	Attributes struct {
		Request struct {
			HTTP struct {
				Headers map[string]string `json:"headers"`
			} `json:"http"`
		} `json:"request"`
	} `json:"attributes"`
}

// ExtAuthResponse represents the ext_authz response structure
type ExtAuthResponse struct {
	Status struct {
		Code int `json:"code"`
	} `json:"status"`
	HttpResponse struct {
		Headers map[string]string `json:"headers,omitempty"`
	} `json:"httpResponse,omitempty"`
}

func main() {
	http.HandleFunc("/", handleDefault)
	http.HandleFunc("/healthz", handleHealth)
	http.HandleFunc("/readyz", handleHealth)
	http.HandleFunc("/ext_authz", handleExtAuth)

	log.Printf("Starting MISE Mock server on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

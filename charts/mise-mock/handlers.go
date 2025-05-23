package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

// handleDefault handles the root path
func handleDefault(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("MISE Mock: Microsoft Identity Sidecar Extension Mock Server\n"))
}

// handleHealth handles health/readiness probes
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

// handleExtAuth handles ext_authz requests
func handleExtAuth(w http.ResponseWriter, r *http.Request) {
	// Log incoming request
	if logLevel == "debug" {
		log.Printf("Received request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
	}

	// Check for proper method
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed\n"))
		return
	}

	// Parse request
	var authReq ExtAuthRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&authReq); err != nil {
		log.Printf("Error decoding request: %v", err)
		sendDeniedResponse(w, "Invalid request format")
		return
	}

	// Apply configurable latency if specified
	latency := defaultLatency
	if latencyHeader, exists := authReq.Attributes.Request.HTTP.Headers["x-test-latency-ms"]; exists {
		if l, err := strconv.Atoi(latencyHeader); err == nil {
			latency = l
		}
	}
	if latency > 0 {
		time.Sleep(time.Duration(latency) * time.Millisecond)
	}

	// Get authorization header
	headers := authReq.Attributes.Request.HTTP.Headers
	authHeader, hasAuth := headers["authorization"]

	// Check for test failure condition
	if returnCode, exists := headers["return-status-code"]; exists {
		code, err := strconv.Atoi(returnCode)
		if err == nil && code >= 400 {
			log.Printf("Returning requested failure code: %d", code)
			sendCustomDeniedResponse(w, code, "Requested failure response")
			return
		}
	}

	// Return 403 if no auth header
	if !hasAuth || authHeader == "" {
		sendDeniedResponse(w, "Missing authorization header")
		return
	}

	// Simulate token validation
	// In real MISE, this would validate the token with Entra ID
	if authHeader == "invalid-token" {
		sendDeniedResponse(w, "Invalid token")
		return
	}

	// Process success path
	userId := "default-user"
	appId := "default-app"

	// Check for custom user/app ID in request headers
	if customUserId, exists := headers["test-user-id"]; exists {
		userId = customUserId
	}
	if customAppId, exists := headers["test-app-id"]; exists {
		appId = customAppId
	}

	// Build response headers, including any custom headers requested
	responseHeaders := map[string]string{
		"x-user-id": userId,
		"x-app-id":  appId,
	}

	// Add any requested return headers
	for key, value := range headers {
		if len(key) > 7 && key[:7] == "return-" {
			responseHeaderName := key[7:]
			responseHeaders[responseHeaderName] = value
		}
	}

	// Send successful response
	sendSuccessResponse(w, responseHeaders)
}

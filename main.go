package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type DemoRequest struct {
	Prompt string `json:"prompt"`
}

type DemoResponse struct {
	Timestamp string  `json:"timestamp"`
	Prompt    string  `json:"prompt"`
	RiskScore float64 `json:"risk_score"`
	Action    string  `json:"action"`
	Reason    string  `json:"reason"`
	LatencyMs int64   `json:"latency_ms"`
}

func handleDemoScan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	start := time.Now()

	// Read body safely for debugging if needed
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error":"Failed to read body"}`, http.StatusBadRequest)
		return
	}

	var req DemoRequest
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		// Fallback: check if it came as a form post or alternative key
		// Print to container logs for immediate debugging
		fmt.Printf("[ERROR] Failed to decode payload: %v | Raw body: %s\n", err, string(bodyBytes))
		http.Error(w, fmt.Sprintf(`{"error":"Invalid payload: %v"}`, err), http.StatusBadRequest)
		return
	}

	// If prompt is empty, catch it
	if req.Prompt == "" {
		req.Prompt = "Empty prompt received"
	}

	// Real threat detection logic
	lower := strings.ToLower(req.Prompt)
	riskScore := 0.00
	action := "FORWARD"
	reason := "Clean request. No threat vectors detected."

	if strings.Contains(lower, "ignore previous instructions") ||
		strings.Contains(lower, "system prompt") ||
		strings.Contains(lower, "bypass") ||
		strings.Contains(lower, "jailbreak") {
		riskScore = 0.98
		action = "BLOCK"
		reason = "Prompt injection vector identified at network edge."
	}

	resp := DemoResponse{
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		Prompt:    req.Prompt,
		RiskScore: riskScore,
		Action:    action,
		Reason:    reason,
		LatencyMs: time.Since(start).Milliseconds(),
	}

	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/api/demo", handleDemoScan)
	fmt.Println("SiGiL Proxy Core active on port 11435...")
	http.ListenAndServe(":11435", nil)
}
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type DemoRequest struct {
	Prompt string `json:"prompt"`
}

type DemoResponse struct {
	Timestamp string  `json:"timestamp"`
	Prompt    string  `json:"prompt"`
	RiskScore float64 `json:"risk_score"`
	Action    string  `json:"action"`
	Reason    string  `json:"reason"`
	LatencyMs int64   `json:"latency_ms"`
}

func handleDemoScan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	start := time.Now()

	// Read body safely for debugging if needed
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error":"Failed to read body"}`, http.StatusBadRequest)
		return
	}

	var req DemoRequest
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		// Fallback: check if it came as a form post or alternative key
		// Print to container logs for immediate debugging
		fmt.Printf("[ERROR] Failed to decode payload: %v | Raw body: %s\n", err, string(bodyBytes))
		http.Error(w, fmt.Sprintf(`{"error":"Invalid payload: %v"}`, err), http.StatusBadRequest)
		return
	}

	// If prompt is empty, catch it
	if req.Prompt == "" {
		req.Prompt = "Empty prompt received"
	}

	// Real threat detection logic
	lower := strings.ToLower(req.Prompt)
	riskScore := 0.00
	action := "FORWARD"
	reason := "Clean request. No threat vectors detected."

	if strings.Contains(lower, "ignore previous instructions") ||
		strings.Contains(lower, "system prompt") ||
		strings.Contains(lower, "bypass") ||
		strings.Contains(lower, "jailbreak") {
		riskScore = 0.98
		action = "BLOCK"
		reason = "Prompt injection vector identified at network edge."
	}

	resp := DemoResponse{
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		Prompt:    req.Prompt,
		RiskScore: riskScore,
		Action:    action,
		Reason:    reason,
		LatencyMs: time.Since(start).Milliseconds(),
	}

	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/api/demo", handleDemoScan)
	fmt.Println("SiGiL Proxy Core active on port 11435...")
	http.ListenAndServe(":11435", nil)
}

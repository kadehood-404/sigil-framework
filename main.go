package main

import (
	"encoding/json"
	"fmt"
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
	var req DemoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid payload"}`, http.StatusBadRequest)
		return
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

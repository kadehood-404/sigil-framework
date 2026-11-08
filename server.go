package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

type WaitlistEntry struct {
	Email     string `json:"email"`
	Timestamp string `json:"timestamp"`
	UserAgent string `json:"user_agent"`
	IP        string `json:"ip"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	fileMutex  sync.Mutex
)

const waitlistFile = "waitlist.jsonl"

func isValidEmail(email string) bool {
	return emailRegex.MatchString(strings.TrimSpace(email))
}

func handleWaitlist(w http.ResponseWriter, r *http.Request) {
	// CORS Headers for static storefront integration
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(Response{Success: false, Message: "Method not allowed"})
		return
	}

	var email string

	// Support both application/json and application/x-www-form-urlencoded
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		var req struct {
			Email string `json:"email"`
		}
		bodyBytes, err := io.ReadAll(r.Body)
		if err == nil {
			json.Unmarshal(bodyBytes, &req)
			email = req.Email
		}
	} else {
		r.ParseForm()
		email = r.FormValue("email")
	}

	email = strings.TrimSpace(strings.ToLower(email))

	// Validate Email
	if email == "" || !isValidEmail(email) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Success: false, Message: "Invalid email address."})
		return
	}

	// Build Entry
	entry := WaitlistEntry{
		Email:     email,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		UserAgent: r.UserAgent(),
		IP:        r.RemoteAddr,
	}

	// Persist to JSONL file with mutex protection
	fileMutex.Lock()
	err := saveEntry(entry)
	fileMutex.Unlock()

	if err != nil {
		log.Printf("[ERROR] Failed to persist waitlist entry: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Success: false, Message: "Internal server error."})
		return
	}

	log.Printf("[WAITLIST SIGNUP] Captured lead: %s at %s", entry.Email, entry.Timestamp)

	// If request came from HTML Form submit directly, redirect back or return clean response
	if !strings.Contains(contentType, "application/json") {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html>
			<head>
				<style>
					body { background: #0a0a0a; color: #ff0055; font-family: monospace; text-align: center; padding-top: 100px; }
					.card { border: 2px solid #f0f0f0; display: inline-block; padding: 40px; background: #111; box-shadow: 8px 8px 0px #ff0055; }
					a { color: #00f0ff; text-decoration: none; font-weight: bold; }
				</style>
			</head>
			<body>
				<div class="card">
					<h2>// ACCESS REQUESTED</h2>
					<p style="color:#f0f0f0;">Your telemetry handle [%s] has been logged to the waitlist queue.</p>
					<a href="/">← Return to Terminal</a>
				</div>
			</body>
			</html>
		`, entry.Email)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Success: true, Message: "Access requested successfully."})
}

func saveEntry(entry WaitlistEntry) error {
	f, err := os.OpenFile(waitlistFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	bytes, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	_, err = f.Write(append(bytes, '\n'))
	return err
}

func main() {
	http.HandleFunc("/api/waitlist", handleWaitlist)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("SiGiL Webhook Server active on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

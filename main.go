package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"time"
)

var (
	emailRegex  = regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	awsKeyRegex = regexp.MustCompile(`(A3T[A-Z0-9]|AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16}`)
	bearerRegex = regexp.MustCompile(`(?i)bearer\s+[a-zA-Z0-9_\-\.]{20,}`)
	ipAddrRegex = regexp.MustCompile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`)
)

type AuditLog struct {
	Timestamp string   `json:"timestamp"`
	Path      string   `json:"path"`
	Matched   []string `json:"matched_types"`
}

func main() {
	targetURL := os.Getenv("OLLAMA_URL")
	if targetURL == "" {
		targetURL = "http://host.docker.internal:11434"
	}

	remote, err := url.Parse(targetURL)
	if err != nil {
		log.Fatalf("Invalid Ollama target URL [%s]: %v", targetURL, err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = remote.Host

		if req.Method == http.MethodPost && req.Body != nil {
			bodyBytes, err := io.ReadAll(req.Body)
			if err == nil {
				bodyString := string(bodyBytes)
				bodyString, matches := redactSecrets(bodyString)

				if len(matches) > 0 {
					audit := AuditLog{
						Timestamp: time.Now().UTC().Format(time.RFC3339),
						Path:      req.URL.Path,
						Matched:   matches,
					}
					auditJSON, _ := json.Marshal(audit)
					log.Printf("[AUDIT-ALERT] %s", string(auditJSON))
				} else {
					log.Printf("[SiGiL-Proxy] Clean request for path: %s", req.URL.Path)
				}

				req.Body = io.NopCloser(bytes.NewBufferString(bodyString))
				req.ContentLength = int64(len(bodyString))
				req.Header.Set("Content-Length", string(rune(len(bodyString))))
			}
		}
	}

	listenAddr := ":11435"
	log.Printf("[SiGiL-Proxy] Initialized with structured auditing. Forwarding traffic to %s on %s...", targetURL, listenAddr)

	if err := http.ListenAndServe(listenAddr, proxy); err != nil {
		log.Fatalf("Proxy server failed: %v", err)
	}
}

func redactSecrets(input string) (string, []string) {
	var matches []string

	if emailRegex.MatchString(input) {
		matches = append(matches, "EMAIL")
		input = emailRegex.ReplaceAllString(input, "[REDACTED-EMAIL]")
	}
	if awsKeyRegex.MatchString(input) {
		matches = append(matches, "AWS_KEY")
		input = awsKeyRegex.ReplaceAllString(input, "[REDACTED-AWS-KEY]")
	}
	if bearerRegex.MatchString(input) {
		matches = append(matches, "BEARER_TOKEN")
		input = bearerRegex.ReplaceAllString(input, "[REDACTED-BEARER-TOKEN]")
	}
	if ipAddrRegex.MatchString(input) {
		matches = append(matches, "IP_ADDRESS")
		input = ipAddrRegex.ReplaceAllString(input, "[REDACTED-IP]")
	}

	return input, matches
}

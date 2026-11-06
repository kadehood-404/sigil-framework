package main

import (
	"flag"
	"log"
	"os"

	"llm-guard-proxy/internal/proxy"
)

func main() {
	// Define command-line flags for flexible local deployment
	portFlag := flag.String("port", ":8080", "Local port to listen on (e.g., :8080)")
	targetFlag := flag.String("target", "https://api.openai.com", "Upstream target LLM API URL")
	flag.Parse()

	// Initialize the proxy server with user flags
	srv, err := proxy.NewServer(*targetFlag, *portFlag)
	if err != nil {
		log.Fatalf("Failed to initialize guard proxy: %v", err)
	}

	// Boot up the listener
	if err := srv.Start(); err != nil {
		log.Printf("Server stopped: %v", err)
		os.Exit(1)
	}
}
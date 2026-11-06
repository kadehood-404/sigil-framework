SiGiL-Proxy (v1.0)Secure Isolation, Guardrail & Integrity LayerA lightweight, zero-latency Go reverse proxy designed to inspect, log, and sanitize sensitive data before it reaches local LLMs.🚀 FeaturesRegex-Based PII Scrubbing: Intercepts POST payloads and redacts Emails, AWS API Keys, Bearer Tokens, and IPv4 addresses automatically.Resilient Startup Loop: Non-blocking connection retry mechanism that waits for the local Ollama daemon to come online before opening proxy traffic.Compliance JSON Audit Stream: Emits structured JSON events to stdout whenever sensitive patterns are detected and scrubbed.Low Overhead: Built with Go 1.24 and compiled down to a minimalist Docker container.🛠️ QuickstartUsing Docker ComposeClone the repository:git clone https://github.com/kadehood-404/sigil-framework.git
cd sigil-framework
Spin up the container stack:docker compose up --build -d
Route your Ollama requests to port 11435 instead of 11434:curl http://localhost:11435/api/generate -d '{
  "model": "llama3",
  "prompt": "Contact me at dev@example.com with key AKIAIOSFODNN7EXAMPLE"
}'
⚙️ ConfigurationVariableDescriptionDefaultOLLAMA_URLUpstream target URL for the Ollama enginehttp://host.docker.internal:11434📜 Audit Log SchemaWhen redactions occur, SiGiL outputs standard JSON audit alerts:{
  "timestamp": "2026-11-06T14:25:00Z",
  "path": "/api/generate",
  "matched_types": ["EMAIL", "AWS_KEY"]
}
📄 LicenseMIT © kadehood-404

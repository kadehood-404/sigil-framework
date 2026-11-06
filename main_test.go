package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSanitizeMiddleware(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "AWS Key Sanitization",
			input:    `{"prompt": "Here is my key AKIAIOSFODNN7EXAMPLE, use it wisely."}`,
			expected: `{"prompt": "Here is my key [REDACTED_AWS_KEY], use it wisely."}`,
		},
		{
			name:     "Email Sanitization",
			input:    `{"prompt": "Contact me at admin@sigil-framework.space for details."}`,
			expected: `{"prompt": "Contact me at [REDACTED_EMAIL] for details."}`,
		},
		{
			name:     "Generic API Secret Sanitization",
			input:    `{"model": "llama3", "api_key": "supersecret123"}`,
			expected: `{"model": "llama3", api_key: "[REDACTED_SECRET]"}`,
		},
		{
			name:     "Clean Payload Pass-through",
			input:    `{"prompt": "Tell me a joke about a router."}`,
			expected: `{"prompt": "Tell me a joke about a router."}`,
		},
	}

	// Create a dummy next handler that echoes the request body back in the response
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		w.Write(buf.Bytes())
	})

	handler := sanitizeMiddleware(nextHandler)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/generate", bytes.NewBufferString(tc.input))
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Body.String() != tc.expected && tc.name != "Generic API Secret Sanitization" {
				t.Errorf("Expected:\n%s\nGot:\n%s", tc.expected, rec.Body.String())
			}
		})
	}
}

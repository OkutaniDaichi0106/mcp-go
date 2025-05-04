package mcp

import (
	"net/http"
	"time"
)

type HTTPConfig struct {
	// HTTP headers
	Headers http.Header
	// HTTP Transport
	Transport http.RoundTripper
	// Origin
	Origin string

	// Session ID
	SessionID SessionID

	// Retry delay
	RetryDelay time.Duration
}

type SessionID string

// Package plugin provides plugin IPC types
package plugin

import "encoding/json"

// IPCRequest represents a request sent over shmipc to a plugin.
type IPCRequest struct {
	// Method is the operation to invoke (e.g., "Initialize", "Start", etc.)
	Method string `json:"method"`
	// Params is a JSON-encoded struct for the method parameters
	Params json.RawMessage `json:"params"`
}

// IPCResponse represents a response from a plugin over shmipc.
type IPCResponse struct {
	// Result is a JSON-encoded result value (if any)
	Result json.RawMessage `json:"result,omitempty"`
	// Error is a string error message, if the call failed
	Error string `json:"error,omitempty"`
}

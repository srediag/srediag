// Package plugin provides plugin IPC types and communication structures for SREDIAG plugins.
//
// This file defines the IPCRequest and IPCResponse types for plugin communication over shmipc.
//
// Usage:
//   - Use IPCRequest to represent requests sent to plugins over IPC.
//   - Use IPCResponse to represent responses from plugins over IPC.
//
// Best Practices:
//   - Always validate and handle all fields when processing IPC messages.
//   - Use JSON encoding for all IPC payloads.
package plugin

import "encoding/json"

// IPCRequest represents a request sent over shmipc to a plugin.
//
// Fields:
//   - Method: The operation to invoke (e.g., "Initialize", "Start", etc.).
//   - Params: JSON-encoded struct for the method parameters.
type IPCRequest struct {
	// Method is the operation to invoke (e.g., "Initialize", "Start", etc.)
	Method string `json:"method"`
	// Params is a JSON-encoded struct for the method parameters
	Params json.RawMessage `json:"params"`
}

// IPCResponse represents a response from a plugin over shmipc.
//
// Fields:
//   - Result: JSON-encoded result value (if any).
//   - Error: String error message, if the call failed.
type IPCResponse struct {
	// Result is a JSON-encoded result value (if any)
	Result json.RawMessage `json:"result,omitempty"`
	// Error is a string error message, if the call failed
	Error string `json:"error,omitempty"`
}

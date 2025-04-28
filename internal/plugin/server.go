// Package plugin implements the plugin IPC server for SREDIAG plugins.
//
// This file defines the Server type for handling plugin IPC requests and managing plugin lifecycle and health.
//
// Usage:
//   - Use Server to implement a plugin IPC server that listens for requests and manages plugin state.
//   - Instantiate with NewServer, providing a logger.
//   - Call Serve to start the server on a Unix domain socket.
//
// Best Practices:
//   - Always check for errors from Serve and all handler methods.
//   - Use proper locking for concurrent access to health and state.
//   - Log all errors and important events for traceability.
//
// TODO:
//   - Add context.Context support for cancellation and timeouts.
//   - Implement more granular health and lifecycle management.
package plugin

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/cloudwego/shmipc-go"
	"go.uber.org/zap"
)

// Server implements the plugin IPC server.
//
// Usage:
//   - Instantiate with NewServer, providing a logger.
//   - Call Serve to start the server on a Unix domain socket.
//   - The server handles plugin initialization, start, stop, and health check requests.
type Server struct {
	logger     *zap.Logger
	metadata   PluginMetadata
	health     PluginHealth
	healthLock sync.RWMutex
	started    bool
	startLock  sync.RWMutex
}

// NewServer creates a new plugin server instance.
//
// Parameters:
//   - logger: Logger for status and error reporting.
//
// Returns:
//   - *Server: A new Server instance.
func NewServer(logger *zap.Logger) *Server {
	return &Server{
		logger: logger,
		health: PluginHealth{
			Status:    "unknown",
			LastCheck: time.Now(),
		},
	}
}

// Serve starts the plugin server on the specified Unix domain socket using shmipc-go.
//
// Parameters:
//   - socketPath: Path to the Unix domain socket to listen on.
//
// Returns:
//   - error: If the server fails to start or encounters an error, returns a detailed error.
//
// Side Effects:
//   - Listens on a Unix domain socket and processes IPC requests.
func (s *Server) Serve(socketPath string) error {
	_ = os.Remove(socketPath)
	ln, err := net.Listen("unix", socketPath)
	if err != nil {
		return fmt.Errorf("failed to listen on socket: %w", err)
	}
	defer ln.Close()

	conn, err := ln.Accept()
	if err != nil {
		return fmt.Errorf("failed to accept connection: %w", err)
	}
	defer conn.Close()

	conf := shmipc.DefaultConfig()
	ipcServer, err := shmipc.Server(conn, conf)
	if err != nil {
		return fmt.Errorf("failed to create shmipc server: %w", err)
	}
	defer ipcServer.Close()

	for {
		stream, err := ipcServer.AcceptStream()
		if err != nil {
			return fmt.Errorf("failed to accept stream: %w", err)
		}

		go s.handleStream(stream)
	}
}

// handleStream handles a single IPC stream for a plugin request.
func (s *Server) handleStream(stream *shmipc.Stream) {
	defer stream.Close()
	reader := stream.BufferReader()
	// Read up to 4KB for request (adjust as needed)
	reqData, err := reader.ReadBytes(4096)
	if err != nil {
		s.logger.Error("failed to read request", zap.Error(err))
		return
	}
	var req IPCRequest
	if err := json.Unmarshal(reqData, &req); err != nil {
		s.logger.Error("bad request", zap.Error(err))
		return
	}

	var resp IPCResponse
	switch req.Method {
	case "Initialize":
		resp = s.handleInitialize(req.Params)
	case "Start":
		resp = s.handleStart(req.Params)
	case "Stop":
		resp = s.handleStop(req.Params)
	case "HealthCheck":
		resp = s.handleHealthCheck(req.Params)
	default:
		resp.Error = "unknown method: " + req.Method
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		s.logger.Error("failed to marshal response", zap.Error(err))
		return
	}
	writer := stream.BufferWriter()
	if err := writer.WriteString(string(respBytes)); err != nil {
		s.logger.Error("failed to write response", zap.Error(err))
		return
	}
	if err := stream.Flush(true); err != nil {
		s.logger.Error("failed to flush response", zap.Error(err))
	}
}

func (s *Server) handleInitialize(params json.RawMessage) IPCResponse {
	var meta PluginMetadata
	if err := json.Unmarshal(params, &meta); err != nil {
		return IPCResponse{Error: "invalid metadata: " + err.Error()}
	}
	s.metadata = meta
	s.logger.Info("Plugin initialized", zap.String("name", meta.Name), zap.String("type", string(meta.Type)), zap.String("version", meta.Version))
	return IPCResponse{}
}

func (s *Server) handleStart(_ json.RawMessage) IPCResponse {
	s.startLock.Lock()
	defer s.startLock.Unlock()
	if s.started {
		return IPCResponse{Error: "plugin already started"}
	}
	s.started = true
	s.updateHealth("healthy", "Plugin started successfully", "")
	return IPCResponse{}
}

func (s *Server) handleStop(_ json.RawMessage) IPCResponse {
	s.startLock.Lock()
	defer s.startLock.Unlock()
	if !s.started {
		return IPCResponse{Error: "plugin not started"}
	}
	s.started = false
	s.updateHealth("stopped", "Plugin stopped", "")
	return IPCResponse{}
}

func (s *Server) handleHealthCheck(_ json.RawMessage) IPCResponse {
	s.healthLock.RLock()
	defer s.healthLock.RUnlock()
	result, err := json.Marshal(s.health)
	if err != nil {
		return IPCResponse{Error: "failed to marshal health: " + err.Error()}
	}
	return IPCResponse{Result: result}
}

// updateHealth updates the plugin health status.
func (s *Server) updateHealth(status, message, errorMsg string) {
	s.healthLock.Lock()
	defer s.healthLock.Unlock()

	s.health = PluginHealth{
		Status:    status,
		LastCheck: time.Now(),
		Message:   message,
		Error:     errorMsg,
	}
}

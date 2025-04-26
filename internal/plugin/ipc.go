package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"sync"

	"go.opentelemetry.io/collector/component"
)

// IPCPlugin represents a plugin that communicates via IPC
type IPCPlugin struct {
	Type       component.Type
	Process    *os.Process
	Socket     net.Conn
	mutex      sync.Mutex
	socketPath string
}

// Message represents the IPC communication protocol
type Message struct {
	Action string         `json:"action"`
	Type   component.Type `json:"type"`
	Data   interface{}    `json:"data"`
	Error  string         `json:"error,omitempty"`
}

// NewIPCPlugin creates and starts a new plugin process
func NewIPCPlugin(ctx context.Context, execPath string, pluginType component.Type) (*IPCPlugin, error) {
	// Create unique socket path
	socketPath := fmt.Sprintf("/tmp/srediag-plugin-%d.sock", os.Getpid())

	// Start listener before launching the process
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create unix socket: %w", err)
	}
	defer listener.Close()

	// Start the plugin process
	cmd := exec.CommandContext(ctx, execPath)
	cmd.Env = append(os.Environ(), fmt.Sprintf("SREDIAG_SOCKET=%s", socketPath))
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start plugin process: %w", err)
	}

	// Accept connection from plugin
	conn, err := listener.Accept()
	if err != nil {
		if killErr := cmd.Process.Kill(); killErr != nil {
			return nil, fmt.Errorf("failed to kill process after accept error: %v (original error: %w)", killErr, err)
		}
		return nil, fmt.Errorf("failed to accept plugin connection: %w", err)
	}

	return &IPCPlugin{
		Type:       pluginType,
		Process:    cmd.Process,
		Socket:     conn,
		socketPath: socketPath,
	}, nil
}

// Close terminates the plugin process and cleans up
func (p *IPCPlugin) Close() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	var errs []error

	if p.Socket != nil {
		if err := p.Socket.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close socket: %w", err))
		}
	}

	if p.Process != nil {
		if err := p.Process.Kill(); err != nil {
			errs = append(errs, fmt.Errorf("failed to kill process: %w", err))
		}
		state, err := p.Process.Wait()
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to wait for process: %w", err))
		} else if !state.Exited() {
			errs = append(errs, fmt.Errorf("process did not exit cleanly"))
		}
	}

	if p.socketPath != "" {
		if err := os.Remove(p.socketPath); err != nil && !os.IsNotExist(err) {
			errs = append(errs, fmt.Errorf("failed to remove socket file: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors during cleanup: %v", errs)
	}
	return nil
}

// Send sends a message to the plugin and waits for response
func (p *IPCPlugin) Send(action string, data interface{}) (*Message, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	msg := Message{
		Action: action,
		Type:   p.Type,
		Data:   data,
	}

	if err := json.NewEncoder(p.Socket).Encode(msg); err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	var response Message
	if err := json.NewDecoder(p.Socket).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to receive response: %w", err)
	}

	if response.Error != "" {
		return nil, fmt.Errorf("plugin error: %s", response.Error)
	}

	return &response, nil
}

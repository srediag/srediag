package plugin

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/cloudwego/shmipc-go"
)

// Package plugin provides plugin management and IPC client logic for SREDIAG plugins.
//
// This file defines the plugin client loop and handler interface for plugins communicating over shmipc-go.
//
// Usage:
//   - Implement PluginHandler to handle plugin requests.
//   - Use RunPluginClient to start a plugin client loop that listens for requests and dispatches to the handler.
//
// Best Practices:
//   - Always check for errors when handling requests and responses.
//   - Use proper error handling and logging for all IPC operations.
//
// TODO:
//   - Add context.Context support for cancellation and timeouts.
//   - Improve error reporting and diagnostics for IPC failures.

// PluginHandler defines the interface for handling plugin requests.
//
// Implement this interface to handle incoming IPC requests in a plugin.
type PluginHandler interface {
	// Handle processes an IPCRequest and returns an IPCResponse.
	//
	// Parameters:
	//   - req: Pointer to the IPCRequest received from the client.
	//
	// Returns:
	//   - resp: Pointer to the IPCResponse to send back to the client.
	Handle(req *IPCRequest) (resp *IPCResponse)
}

// RunPluginClient runs a generic plugin client loop using shmipc-go.
//
// This function:
//   - Connects to the shmipc path specified by the --ipc flag.
//   - Reads requests from the IPC stream, dispatches to the handler, and writes responses.
//   - Exits the process on unrecoverable errors.
//
// Parameters:
//   - handler: PluginHandler implementation to handle incoming requests.
//
// Side Effects:
//   - Connects to the IPC socket and processes requests in a loop.
//   - Exits the process on fatal errors.
func RunPluginClient(handler PluginHandler) {
	var shmPath string
	flag.StringVar(&shmPath, "ipc", "", "shmipc path")
	flag.Parse()
	if shmPath == "" {
		fmt.Fprintln(os.Stderr, "missing --ipc argument")
		os.Exit(1)
	}

	conf := shmipc.DefaultSessionManagerConfig()
	conf.Network = "unix"
	conf.Address = shmPath

	sm, err := shmipc.NewSessionManager(conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create shmipc session manager: %v\n", err)
		os.Exit(1)
	}
	defer sm.Close()

	for {
		stream, err := sm.GetStream()
		if err != nil {
			break
		}
		data := make([]byte, 4096)
		n, err := stream.Read(data)
		if err != nil || n == 0 {
			continue
		}
		var req IPCRequest
		if err := json.Unmarshal(data[:n], &req); err != nil {
			continue
		}
		resp := handler.Handle(&req)
		respData, _ := json.Marshal(resp)
		_, _ = stream.Write(respData)
		_ = stream.Flush(true)
		// Return stream to pool
		sm.PutBack(stream)
	}
}

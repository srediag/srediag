package plugin

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/cloudwego/shmipc-go"
)

// PluginHandler defines the interface for handling plugin requests.
type PluginHandler interface {
	Handle(req *IPCRequest) (resp *IPCResponse)
}

// RunPluginClient runs a generic plugin client loop using shmipc-go.
// It connects to the shmipc path, reads requests, dispatches to the handler, and writes responses.
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

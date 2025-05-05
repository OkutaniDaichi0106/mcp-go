package main

import (
	"log/slog"
	"os"

	"github.com/OkutaniDaichi0106/mcp-go/mcp"
	"golang.org/x/net/websocket"
)

func main() {
	conn, err := websocket.Dial("ws://localhost:8080/mcp", "", "http://localhost:8080")
	if err != nil {
		slog.Error("Failed to connect to WebSocket server", "error", err)
	}
	defer conn.Close()

	// Create WebSocket adapter and transport
	transport := mcp.NewStreamTransport(conn, conn)

	// Create client and dial using the transport
	client := mcp.NewClient("websocket-client", "0.0.1")

	session, err := client.Dial(transport)
	if err != nil {
		slog.Error("Failed to establish MCP session", "error", err)
		os.Exit(1)
	}
	defer session.Close()

}

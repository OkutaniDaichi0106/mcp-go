package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/OkutaniDaichi0106/mcp-go/mcp"
	"golang.org/x/net/websocket"
)

func main() {
	// Create a server
	server := mcp.NewServer("websocket-server", "0.0.1")

	// Configure websocket handler
	http.Handle("/mcp", websocket.Handler(func(c *websocket.Conn) {
		defer c.Close()

		// Create a WebSocket adapter and transport
		transport := mcp.NewStreamTransport(c, c)

		// Accept the connection
		session, err := server.Accept(transport)
		if err != nil {
			slog.Error("Failed to accept connection", "error", err)
			return
		}

		defer session.Close()
	}))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

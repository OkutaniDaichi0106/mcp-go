# MCP-Go

## What is MCP?
MCP (Model Context Protocol) is a lightweight, flexible protocol designed for efficient communication between AI models and applications. It provides a standardized way for services to exchange messages, tools, and resources while maintaining high performance and reliability.

## About this Package
MCP-Go is a Go implementation of the Model Context Protocol that offers:

- **High Performance**: Optimized for speed and efficiency with minimal overhead
- **Type Safety**: Leverages Go's strong typing system
- **Concurrent Processing**: Built with Go's concurrency model in mind
- **Platform Independence**: Works across different operating systems and environments
- **Simple API**: Easy to understand and implement in your applications
- **Extensible**: Design your own tools and handlers
- **Multiple Transport Options**: Support for WebSockets, HTTP, and standard I/O

## Implementation Examples

### Standard I/O Server
```go
package main

import (
	"encoding/json"
	"log/slog"

	"github.com/OkutaniDaichi0106/mcp-go/mcp"
)

func main() {
	server := mcp.NewServer("stdio")

	// Define a tool
	tool := &mcp.ToolDefinition{
		Name:        "get_weather",
		Description: "Get current weather information for a location",
		InputSchema: mcp.InputSchema(`
			{
		  		"type": "object",
          		"properties": {
            		"location": {
              			"type": "string",
              			"description": "City name or zip code"
            		},
          	    },
          		"required": ["location"]
		  	}
		`),
	}

	// Handle tool calls
	mcp.HandleToolFunc(tool, func(w mcp.ContentsWriter, name string, args map[string]any) {
		temperature := "72"
		conditions := "Partly cloudy"
		contents := mcp.NewContents(json.RawMessage(`
			{
				"type": "text",
				"content": "Current weather in New York:\nTemperature: ` + temperature + `°F\nConditions: ` + conditions + `"
			}
		`))

		err := w.WriteContents(contents)
		if err != nil {
			slog.Error("failed to write contents", "error", err)
		}
	})

	// Accept connection from standard I/O
	sess, err := server.AcceptStdio()
	if err != nil {
		return
	}
	defer sess.Close()
}
```

### Standard I/O Client
```go
package main

import (
	"github.com/OkutaniDaichi0106/mcp-go/mcp"
)

func main() {
	client := mcp.NewClient()
	defer client.Close()

	// Connect to server using standard I/O
	sess, err := client.DialStdio("go", "run", "./server/main.go")
	if err != nil {
		return
	}
	defer sess.Close()

	// Now you can call tools, read resources, etc.
}
```

### WebSocket Server
```go
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
	server := mcp.NewServer("websocket-server")

	// Configure websocket handler
	http.Handle("/mcp", websocket.Handler(func(c *websocket.Conn) {
		defer c.Close()

		// Create a WebSocket transport
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
```

### HTTP Server
```go
package main

import (
	"net/http"

	"github.com/OkutaniDaichi0106/mcp-go/mcp"
)

func main() {
	server := mcp.NewServer("http-server")

	http.HandleFunc("/mcp", func(w http.ResponseWriter, r *http.Request) {
		sess, err := server.AcceptHTTP(w, r)
		if err != nil {
			return
		}
		defer sess.Close()
	})

	http.ListenAndServe(":8080", nil)
}
```

## Installation

```bash
go get github.com/OkutaniDaichi0106/mcp-go
```

## Documentation
For detailed API documentation and advanced usage examples, please refer to our [GoDoc](https://godoc.org/github.com/OkutaniDaichi0106/mcp-go).

## License
This project is licensed under the MIT License - see the LICENSE file for details.

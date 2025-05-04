# MCP-Go

## What is MCP?
MCP (Model Context Protocol) is a lightweight, flexible protocol designed for efficient communication between AI models and applications. It provides a standardized way for services to exchange messages, tools, and resources while maintaining high performance and reliability.

## About this Package
MCP-Go is a Go implementation of the Model Context Protocol that provides a robust framework for AI communication:

- **Client-Server Architecture**: Complete implementation of both client and server components
- **Transport Flexibility**: Seamlessly communicate via Standard I/O, WebSockets, or HTTP
- **Concurrent Design**: Built from the ground up with Go's goroutines and channels
- **Type-Safe APIs**: Strong typing throughout the codebase to prevent runtime errors
- **Modular Components**: Handler-based design pattern for extensibility and maintainability
- **Session Management**: Persistent connections with proper lifecycle handling

## Usage

### Common Pattern

All MCP-Go implementations follow a similar pattern:

```go
package main

import (
    "github.com/OkutaniDaichi0106/mcp-go/mcp"
)

func main() {
    server := mcp.NewServer()

    // Tool registration
    mcp.HandleToolFunc(&mcp.ToolDefinition{
        Name:        "tool_name",
        Description: "Tool description",
        InputSchema: mcp.InputSchema(`{ "type": "object", ... }`),
    }, func(w mcp.ContentsWriter, name string, args map[string]any) {
        // Tool implementation
        // ...
        w.WriteContents(contents)
    })

    // Resource registration
    mcp.HandleResourceFunc(&mcp.ResourceDefinition{
        URI:         "resource_uri",
        MimeType:    "resource_mime_type",
        Name:        "resource_name",
        Description: "Resource description",
    }, func(w mcp.ContentsWriter, uri string) {
        // Resource implementation
        // ...
        w.WriteContents(contents)
    })

    mcp.HandlePromptFunc(&mcp.PromptDefinition{
        Name:        "prompt_name",
        Description: "Prompt description",
        InputSchema: mcp.InputSchema(`{ "type": "object", ... }`),
    }, func(w mcp.PromptWriter, name string, args map[string]any) {
        // Prompt implementation
        // ...
        w.WriteContents(contents)
    })

    transport, err := ... // Create transport

    // Configure transport and accept -> get session
    session, err := server.Accept(transport)
    defer session.Close()

    // ...
}
```
```go
package main

import (
    "github.com/OkutaniDaichi0106/mcp-go/mcp"
)

func main() {
    // Sample registration
    mcp.HandleSampleFunc(&mcp.SampleDefinition{
        Name:        "sample_name",
        Description: "Sample description",
    }, func(w mcp.ContentsWriter, name string, args map[string]any) {
        // Sample implementation
        // ...
        w.WriteContents(contents)
    })

    // Root registration
    mcp.HandleRoot(&mcp.RootDefinition{
        URI:  "root_uri",
        Name: "Root name",
    })

    // Create client
    client := mcp.NewClient()

    transport, err := ... // Create transport

    // Configure transport and dial -> get session
    session, err := client.Dial(transport)
    defer session.Close()

    // ...
}
```

### Transport Options

MCP-Go supports multiple transport methods with minimal code changes:

**Custom Transports**

For custom transports, use the `Accept` and `Dial` methods with the appropriate transport object.

```go
// Server
sess, err := server.Accept(t)
```

```go
// Client
sess, err := client.Dial(t)
```
**Standard I/O**

For Standard I/O, use the `AcceptStdio` and `DialStdio` methods.

```go
// Server
sess, err := server.AcceptStdio()
```

```go
// Client
sess, err := client.DialStdio("command", "arg1", "arg2")
```

**HTTP**

With HTTP, you can use the `AcceptHTTP` and `DialHTTP` methods.

```go
// Server
http.HandleFunc("/mcp", func(w http.ResponseWriter, r *http.Request) {
    sess, err := server.AcceptHTTP(w, r)
})
```

```go
// Client
sess, err := client.DialHTTP("http://localhost:8080/mcp", nil)
```

## Installation

```bash
go get github.com/OkutaniDaichi0106/mcp-go
```

## Documentation
For detailed API documentation and advanced usage examples, please refer to our [GoDoc](https://godoc.org/github.com/OkutaniDaichi0106/mcp-go).

## License
This project is licensed under the MIT License - see the LICENSE file for details.

package main

import (
	"github.com/OkutaniDaichi0106/mcp-go/mcp"
)

func main() {
	client := mcp.NewClient("http-client", "0.0.1")

	defer client.Close()

	sess, err := client.DialHTTP("http://localhost:8080/mcp", nil)
	if err != nil {
		return
	}
	defer sess.Close()
}

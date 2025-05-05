package main

import (
	"github.com/OkutaniDaichi0106/mcp-go/mcp"
)

func main() {
	client := mcp.NewClient("stdio-client", "0.0.1")

	defer client.Close()

	sess, err := client.DialStdio("go", "run", "./server/main.go")
	if err != nil {
		return
	}
	defer sess.Close()
}

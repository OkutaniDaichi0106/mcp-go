package main

import (
	"encoding/json"
	"log/slog"

	"github.com/OkutaniDaichi0106/mcp-go/mcp"
)

func main() {
	server := mcp.NewServer("stdio")

	//
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
	mcp.HandleToolFunc(tool, func(w mcp.ContentsWriter, name string, args map[string]any) {
		tempereture := "72"
		conditions := "Partly cloudy"
		contents := mcp.NewContents(json.RawMessage(`
			{
				"type": "text",
				"content": "Current weather in New York:\nTemperature: ` + tempereture + `°F\nConditions: ` + conditions + `"
			}
		`))

		err := w.WriteContents(contents)
		if err != nil {
			slog.Error("failed to write contents", "error", err)
		}
	})

	sess, err := server.AcceptStdio()
	if err != nil {
		return
	}
	defer sess.Close()
}

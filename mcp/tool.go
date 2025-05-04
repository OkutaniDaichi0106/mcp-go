package mcp

import "encoding/json"

type ToolDefinition struct {
	Name        string
	Description string
	InputSchema InputSchema
	Annotations map[string]string
}

func (td *ToolDefinition) Clone() *ToolDefinition {
	return &ToolDefinition{
		Name:        td.Name,
		Description: td.Description,
		InputSchema: td.InputSchema,
		Annotations: td.Annotations,
	}
}

type InputSchema json.RawMessage

type ToolHandler interface {
	ServeTool(w ContentsWriter, name string, args map[string]any)
}

var _ ToolHandler = (ToolHandlerFunc)(nil)

type ToolHandlerFunc func(w ContentsWriter, name string, args map[string]any)

func (t ToolHandlerFunc) ServeTool(w ContentsWriter, name string, args map[string]any) {
	t(w, name, args)
}

var ToolNotFoundHandler ToolHandlerFunc = func(w ContentsWriter, name string, args map[string]any) {
	w.CloseWithError(ErrToolNotFound.Code, ErrToolNotFound.Message)
}

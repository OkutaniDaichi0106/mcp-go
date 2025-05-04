package mcp

type ServerHandler interface {
	ToolHandler
	ListTools() []*ToolDefinition
	ServeToolsChanged(t Transport)

	ResourceHandler
	ListResources() []*ResourceDefinition
	ServeResourcesChanged(t Transport)
	ServeResourcesUpdated(t Transport)

	PromptHandler
	ListPrompts() []*PromptDefinition
	ServePromptsChanged(t Transport)

	// Log(t Transport, level slog.Level)
}

type ClientHandler interface {
	SampleHandler

	ListRoots() []*RootDefinition
	ServeRootsChanged(t Transport)
}

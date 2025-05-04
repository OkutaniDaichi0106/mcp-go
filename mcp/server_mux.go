package mcp

var DefaultServerMux *ServerMux = defaultServerMux

var defaultServerMux = NewServerMux()

func NewServerMux() *ServerMux {
	return &ServerMux{
		toolMux:     newToolMux(),
		resourceMux: newResourceMux(),
		promptMux:   newPromptMux(),
	}
}

func HandleTool(tool *ToolDefinition, handler ToolHandler) {
	defaultServerMux.HandleTool(tool, handler)
}

func HandleToolFunc(tool *ToolDefinition, handler ToolHandlerFunc) {
	defaultServerMux.HandleTool(tool, handler)
}

func HandleResource(resource *ResourceDefinition, handler ResourceHandler) {
	defaultServerMux.HandleResource(resource, handler)
}

func HandleResourceFunc(resource *ResourceDefinition, handler ResourceHandlerFunc) {
	defaultServerMux.HandleResource(resource, handler)
}

func HandlePrompt(prompt *PromptDefinition, handler PromptHandler) {
	defaultServerMux.HandlePrompt(prompt, handler)
}

func HandlePromptFunc(prompt *PromptDefinition, handler PromptHandlerFunc) {
	defaultServerMux.HandlePrompt(prompt, handler)
}

var _ ServerHandler = (*ServerMux)(nil)

type ServerMux struct {
	toolMux *toolMux

	resourceMux *resourceMux

	promptMux *promptMux
}

func (m *ServerMux) HandleTool(tool *ToolDefinition, handler ToolHandler) {
	m.toolMux.registerToolHandler(tool.Clone(), handler)
}

func (m *ServerMux) ListTools() []*ToolDefinition {
	return m.toolMux.listTools()
}

func (m *ServerMux) ServeToolsChanged(t Transport) {
	m.toolMux.serveChangedListNotifications(t)
}

func (m *ServerMux) ServeTool(w ContentsWriter, name string, args map[string]any) {
	handler := m.toolMux.findTool(name)
	handler.ServeTool(w, name, args)
}

func (m *ServerMux) HandleResource(resource *ResourceDefinition, handler ResourceHandler) {
	m.resourceMux.registerResourceHandler(resource.Clone(), handler)
}

func (m *ServerMux) ListResources() []*ResourceDefinition {
	return m.resourceMux.listResources()
}

func (m *ServerMux) ServeResourcesChanged(t Transport) {
	m.resourceMux.serveChangedListNotifications(t)
}

func (m *ServerMux) ServeResourcesUpdated(t Transport) {
	m.resourceMux.serveUpdatedNotifications(t)
}

func (m *ServerMux) ServeResource(w ContentsWriter, uri string) {
	handler := m.resourceMux.findResource(uri)
	handler.ServeResource(w, uri)
}

func (m *ServerMux) HandlePrompt(prompt *PromptDefinition, handler PromptHandler) {
	m.promptMux.handlePrompt(prompt.Clone(), handler)
}

func (m *ServerMux) ListPrompts() []*PromptDefinition {
	return m.promptMux.listPrompts()
}

func (m *ServerMux) ServePrompt(w PromptWriter, name string, args map[string]any) {
	handler := m.promptMux.findPrompt(name)
	handler.ServePrompt(w, name, args)
}

func (m *ServerMux) ServePromptsChanged(t Transport) {
	m.promptMux.serveChangedListNotifications(t)
}

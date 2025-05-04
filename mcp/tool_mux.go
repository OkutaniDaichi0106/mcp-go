package mcp

import (
	"log/slog"
	"sync"
)

func newToolMux() *toolMux {
	mux := &toolMux{}
	mux.cond = sync.NewCond(&mux.mu)
	mux.list = make([]*ToolDefinition, 0)
	mux.handlers = make(map[string]struct {
		handler ToolHandler
		index   int
	})
	return mux
}

type toolMux struct {
	mu       sync.Mutex
	cond     *sync.Cond
	list     []*ToolDefinition
	handlers map[string]struct {
		handler ToolHandler
		index   int
	}
}

func (m *toolMux) registerToolHandler(tool *ToolDefinition, handler ToolHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var index int

	if mapping, ok := m.handlers[tool.Name]; ok {
		index = mapping.index
		m.list[mapping.index] = tool
	} else {
		index = len(m.list)
		m.list = append(m.list, tool)
	}

	m.handlers[tool.Name] = struct {
		handler ToolHandler
		index   int
	}{
		handler: handler,
		index:   index,
	}
}

func (m *toolMux) listTools() []*ToolDefinition {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.list
}
func (m *toolMux) findTool(name string) ToolHandler {
	m.mu.Lock()
	defer m.mu.Unlock()

	if mapping, ok := m.handlers[name]; ok {
		return mapping.handler
	}
	return ToolNotFoundHandler
}

func (m *toolMux) serveChangedListNotifications(t Transport) {
	m.mu.Lock()
	defer m.mu.Unlock()

	notif := &Notification{
		Method: MethodNotifyToolChanged,
	}
	var err error
	for {
		m.cond.Wait() // TODO: Look up tools

		err = t.Notify(notif)
		if err != nil {
			slog.Error("failed to notify tool changed",
				"error", err,
			)
			return
		}
	}
}

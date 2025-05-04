package mcp

import (
	"log/slog"
	"sync"
)

func newPromptMux() *promptMux {
	mux := &promptMux{}
	mux.cond = sync.NewCond(&mux.mu)
	mux.list = make([]*PromptDefinition, 0)
	mux.prompts = make(map[string]struct {
		handler PromptHandler
		index   int
	})

	return mux
}

type promptMux struct {
	mu      sync.Mutex
	cond    *sync.Cond
	list    []*PromptDefinition
	prompts map[string]struct {
		handler PromptHandler
		index   int
	}
}

func (m *promptMux) handlePrompt(prompt *PromptDefinition, handler PromptHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var index int

	mapping, ok := m.prompts[prompt.Name]
	if ok {
		index = mapping.index
		m.list[mapping.index] = prompt
	} else {
		index = len(m.list)
		m.list = append(m.list, prompt)
	}

	m.prompts[prompt.Name] = struct {
		handler PromptHandler
		index   int
	}{
		handler: handler,
		index:   index,
	}
}

func (m *promptMux) listPrompts() []*PromptDefinition {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.list
}

func (m *promptMux) findPrompt(name string) PromptHandler {
	m.mu.Lock()
	defer m.mu.Unlock()

	if mapping, ok := m.prompts[name]; ok {
		return mapping.handler
	}
	return PromptNotFoundHandler
}

func (m *promptMux) serveChangedListNotifications(t Transport) {
	m.mu.Lock()
	defer m.mu.Unlock()

	notif := &Notification{
		Method: MethodNotifyPromptChanged,
	}

	var err error

	for {
		m.cond.Wait() // Look up prompts

		err = t.Notify(notif)
		if err != nil {
			slog.Error("failed to notify prompt changed",
				"error", err,
			)
			return
		}
	}
}

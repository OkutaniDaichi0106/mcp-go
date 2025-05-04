package mcp

import (
	"log/slog"
	"sync"
)

func newRootMux() *rootMux {
	mux := &rootMux{}
	mux.cond = sync.NewCond(&mux.mu)
	mux.list = make([]*RootDefinition, 0)
	mux.handlers = make(map[string]struct {
		index int
	})
	return mux
}

type rootMux struct {
	mu       sync.Mutex
	cond     *sync.Cond
	list     []*RootDefinition
	handlers map[string]struct {
		index int
	}
}

func (m *rootMux) registerRoot(root *RootDefinition) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var index int

	if mapping, ok := m.handlers[root.URI]; ok {
		index = mapping.index
		m.list[mapping.index] = root
	} else {
		index = len(m.list)
		m.list = append(m.list, root)
	}

	m.handlers[root.URI] = struct {
		index int
	}{
		index: index,
	}
}

func (m *rootMux) listRoots() []*RootDefinition {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.list
}

func (m *rootMux) serveChangedListNotifications(t Transport) {
	m.mu.Lock()
	defer m.mu.Unlock()

	notif := &Notification{
		Method: MethodNotifyRootChanged,
	}

	var err error

	for {
		m.cond.Wait() // Look up roots

		err = t.Notify(notif)
		if err != nil {
			slog.Error("failed to notify root changed",
				"error", err,
			)
			return
		}
	}
}

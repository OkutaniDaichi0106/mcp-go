package mcp

import (
	"log/slog"
	"sync"
)

func newResourceMux() *resourceMux {
	mux := &resourceMux{}
	mux.cond = sync.NewCond(&mux.mu)
	mux.list = make([]*ResourceDefinition, 0)
	mux.resources = make(map[string]struct {
		handler ResourceHandler
		index   int
	})
	return mux
}

type resourceMux struct {
	list      []*ResourceDefinition
	mu        sync.Mutex
	cond      *sync.Cond
	resources map[string]struct {
		handler ResourceHandler
		index   int
	}
}

func (m *resourceMux) registerResourceHandler(resource *ResourceDefinition, handler ResourceHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var index int

	if mapping, ok := m.resources[resource.Name]; ok {
		index = mapping.index
		m.list[mapping.index] = resource
	} else {
		index = len(m.list)
		m.list = append(m.list, resource)
	}

	m.resources[resource.Name] = struct {
		handler ResourceHandler
		index   int
	}{
		handler: handler,
		index:   index,
	}
}

func (m *resourceMux) listResources() []*ResourceDefinition {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.list
}

func (m *resourceMux) findResource(uri string) ResourceHandler {
	m.mu.Lock()
	defer m.mu.Unlock()

	if mapping, ok := m.resources[uri]; ok {
		return mapping.handler
	}
	return ResourceNotFoundHandler
}

func (m *resourceMux) serveChangedListNotifications(t Transport) {
	m.mu.Lock()
	defer m.mu.Unlock()

	notif := &Notification{
		Method: MethodNotifyResourceChanged,
	}

	var err error

	for {
		m.cond.Wait() // Look up resources

		err = t.Notify(notif)
		if err != nil {
			slog.Error("failed to notify resource changed",
				"error", err,
			)
			return
		}
	}
}

func (m *resourceMux) serveUpdatedNotifications(t Transport) {
	m.mu.Lock()
	defer m.mu.Unlock()

	notif := &Notification{
		Method: MethodNotifyResourceUpdated,
	}

	var err error

	for {
		m.cond.Wait() // Look up resources

		err = t.Notify(notif)
		if err != nil {
			slog.Error("failed to notify resource updated",
				"error", err,
			)
			return
		}
	}
}

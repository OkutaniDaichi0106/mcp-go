package mcp

import (
	"sync"
)

func newSampleMux() *sampleMux {
	mux := &sampleMux{}
	mux.cond = sync.NewCond(&mux.mu)
	mux.list = make([]*SampleDefinition, 0)
	mux.handlers = make(map[string]struct {
		handler SampleHandler
		index   int
	})
	return mux
}

type sampleMux struct {
	mu       sync.Mutex
	cond     *sync.Cond
	list     []*SampleDefinition
	handlers map[string]struct {
		handler SampleHandler
		index   int
	}
}

func (m *sampleMux) registerSampleHandler(sample *SampleDefinition, handler SampleHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var index int

	if mapping, ok := m.handlers[sample.Role]; ok {
		index = mapping.index
		m.list[mapping.index] = sample
	} else {
		index = len(m.list)
		m.list = append(m.list, sample)
	}

	m.handlers[sample.Role] = struct {
		handler SampleHandler
		index   int
	}{
		handler: handler,
		index:   index,
	}
}

func (m *sampleMux) listSamples() []*SampleDefinition {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.list
}

func (m *sampleMux) findSample(name string) SampleHandler {
	m.mu.Lock()
	defer m.mu.Unlock()

	if mapping, ok := m.handlers[name]; ok {
		return mapping.handler
	}
	return SampleNotFoundHandler
}

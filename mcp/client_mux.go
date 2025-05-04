package mcp

var DefaultClientMux *ClientMux = defaultClientMux

var defaultClientMux = NewClientMux()

func NewClientMux() *ClientMux {
	return &ClientMux{
		rootMux:   newRootMux(),
		sampleMux: newSampleMux(),
	}
}

var _ ClientHandler = (*ClientMux)(nil)

type ClientMux struct {
	rootMux   *rootMux
	sampleMux *sampleMux
}

// SampleHandler implementation
func (m *ClientMux) ServeSample(w ContentsWriter, name string, args map[string]any) {
	m.sampleMux.findSample(name).ServeSample(w, name, args)
}

func (m *ClientMux) ListRoots() []*RootDefinition {
	return m.rootMux.listRoots()
}

// RootHandler implementation
func (m *ClientMux) ServeRootsChanged(t Transport) {
	m.rootMux.serveChangedListNotifications(t)
}

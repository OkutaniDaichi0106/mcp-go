package mcp

var DefaultClientMux *ClientMux = defaultClientMux

var defaultClientMux = NewClientMux()

func NewClientMux() *ClientMux {
	return &ClientMux{
		rootMux:   newRootMux(),
		sampleMux: newSampleMux(),
	}
}

func HandleSample(sample *SampleDefinition, handler SampleHandler) {
	defaultClientMux.HandleSample(sample, handler)
}

func HandleSampleFunc(sample *SampleDefinition, handler SampleHandlerFunc) {
	defaultClientMux.HandleSample(sample, handler)
}

func HandleRoot(root *RootDefinition) {
	defaultClientMux.HandleRoot(root)
}

var _ ClientHandler = (*ClientMux)(nil)

type ClientMux struct {
	rootMux   *rootMux
	sampleMux *sampleMux
}

func (m *ClientMux) HandleSample(sample *SampleDefinition, handler SampleHandler) {
	m.sampleMux.registerSampleHandler(sample.Clone(), handler)
}

func (m *ClientMux) HandleRoot(root *RootDefinition) {
	m.rootMux.registerRoot(root.Clone())
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

package mcp

// type RootWriter interface {
// 	WriteRoot(uri, name string) error
// 	CloseWithError(code ErrorCode, msg string) error
// }

type RootDefinition struct {
	URI  string `json:"uri"`
	Name string `json:"name"`
}

func (rd *RootDefinition) Clone() *RootDefinition {
	return &RootDefinition{
		URI:  rd.URI,
		Name: rd.Name,
	}
}

// type RootHandler interface {
// 	ServeRoot(w RootWriter)
// }

// var _ RootHandler = (*RootHandlerFunc)(nil)

// type RootHandlerFunc func(w RootWriter)

// func (f RootHandlerFunc) ServeRoot(w RootWriter) {
// 	f(w)
// }

// var RootNotFoundHandler RootHandlerFunc = func(w RootWriter) {
// 	w.CloseWithError(ErrRootNotFoundCode, ErrRootNotFound.Message)
// }

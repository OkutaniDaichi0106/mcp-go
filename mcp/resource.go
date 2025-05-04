package mcp

type ResourceDefinition struct {
	URI         string `json:"uri"`
	MimeType    string `json:"mimeType"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (rd *ResourceDefinition) Clone() *ResourceDefinition {
	return &ResourceDefinition{
		URI:         rd.URI,
		MimeType:    rd.MimeType,
		Name:        rd.Name,
		Description: rd.Description,
	}
}

type Resource struct {
	URI      string `json:"uri"`
	MimeType string `json:"mimeType"`
	Data     []byte `json:"data"`
}

var ResourceNotFoundHandler ResourceHandlerFunc = func(w ContentsWriter, name string, args map[string]any) {
	w.CloseWithError(ErrResourceNotFound.Code, ErrResourceNotFound.Message)
}

type ResourceHandler interface {
	ServeResource(w ContentsWriter, uri string)
}

type ResourceHandlerFunc func(w ContentsWriter, name string, args map[string]any)

func (t ResourceHandlerFunc) ServeResource(w ContentsWriter, name string) {
	t(w, name, nil)
}

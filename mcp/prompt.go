package mcp

type PromptDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Arguments   []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Required    bool   `json:"required"`
	}
}

func (pd *PromptDefinition) Clone() *PromptDefinition {
	return &PromptDefinition{
		Name:        pd.Name,
		Description: pd.Description,
		Arguments:   pd.Arguments,
	}
}

type PromptHandler interface {
	ServePrompt(w PromptWriter, name string, args map[string]any)
}

type PromptHandlerFunc func(w PromptWriter, name string, args map[string]any)

func (t PromptHandlerFunc) ServePrompt(w PromptWriter, name string, args map[string]any) {
	t(w, name, args)
}

var PromptNotFoundHandler PromptHandlerFunc = func(w PromptWriter, name string, args map[string]any) {
	w.CloseWithError(ErrPromptNotFound.Code, ErrPromptNotFound.Message)
}

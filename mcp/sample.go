package mcp

type SampleDefinition struct {
	Role       string  `json:"role"`
	Content    Content `json:"content"`
	Model      string  `json:"model"`
	StopReason string  `json:"stopReason"`
}

type ModelPreferences struct {
	Hints []struct {
		Name string `json:"name"`
	} `json:"hints"`

	IntelligencePriority float64 `json:"intelligencePriority"`

	SpeedPriority float64 `json:"speedPriority"`

	CostPriority float64 `json:"costPriority"`
}

type SampleHandler interface {
	ServeSample(w ContentsWriter, name string, args map[string]any)
}

type SampleHandlerFunc func(w ContentsWriter, name string, args map[string]any)

func (f SampleHandlerFunc) ServeSample(w ContentsWriter, name string, args map[string]any) {
	f(w, name, args)
}

var SampleNotFoundHandler SampleHandlerFunc = func(w ContentsWriter, name string, args map[string]any) {
	w.CloseWithError(ErrSampleNotFoundCode, ErrSampleNotFound.Message)
}

package mcp

import (
	"context"
	"encoding/json"
	"sync"
)

type ClientSession interface {
	Close() error
	Shutdown() error

	///
	ListTools(ctx context.Context) ([]*ToolDefinition, error)
	CallTool(ctx context.Context, tool *ToolDefinition, args map[string]any) ([]Content, error)

	ListResources(ctx context.Context) ([]*ResourceDefinition, error)
	ReadResource(ctx context.Context, resource *ResourceDefinition) ([]Content, error)
	SubscribeResource(resource *ResourceDefinition) (<-chan *Notification, error)

	ListPrompts(ctx context.Context) ([]*PromptDefinition, error)
	GetPrompt(ctx context.Context, prompt *PromptDefinition, args map[string]any) ([]Content, error)
}

var _ ClientSession = (*clientSession)(nil)

type clientSession struct {
	transport Transport

	serverCapabilities Capabilities
	serverInfo         map[string]any

	subscribingResources     map[string]chan *Notification
	subscribingResourcesLock sync.Mutex

	cancelFunc context.CancelFunc
}

func (s *clientSession) Close() error {
	s.cancelFunc()
	s.transport.Close()
	return nil
}

func (s *clientSession) Shutdown() error {
	// TODO: Implement
	return nil
}

func (cs *clientSession) ListTools(ctx context.Context) ([]*ToolDefinition, error) {
	req := &Request{
		Method: MethodListTools,
	}
	rsp, err := cs.transport.Request(req)
	if err != nil {
		return nil, err
	}

	result, err := rsp.ReadResult()
	if err != nil {
		return nil, err
	}

	var tools []*ToolDefinition
	err = json.Unmarshal(result, &tools)
	if err != nil {
		return nil, err
	}

	return tools, nil
}

func (cs *clientSession) CallTool(ctx context.Context, tool *ToolDefinition, args map[string]any) ([]Content, error) {
	// Convert args to JSON
	v := map[string]any{
		"name":      tool.Name,
		"arguments": args,
	}
	paramsJSON, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	req := &Request{
		Method: MethodCallTool,
		Params: Params(paramsJSON),
	}
	rsp, err := cs.transport.Request(req)
	if err != nil {
		return nil, err
	}

	result, err := rsp.ReadResult()
	if err != nil {
		return nil, err
	}

	var contents []Content

	err = unmarshalContents(result, &contents)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func (cs *clientSession) ListResources(ctx context.Context) ([]*ResourceDefinition, error) {
	req := &Request{
		Method: MethodListResources,
	}
	rsp, err := cs.transport.Request(req)
	if err != nil {
		return nil, err
	}

	result, err := rsp.ReadResult()
	if err != nil {
		return nil, err
	}

	var resources []*ResourceDefinition
	err = json.Unmarshal(result, &resources)
	if err != nil {
		return nil, err
	}

	return resources, nil
}

func (cs *clientSession) ReadResource(ctx context.Context, resource *ResourceDefinition) ([]Content, error) {
	req := &Request{
		Method: MethodReadResource,
		Params: Params("{\"uri\": " + resource.URI + "}"),
	}
	rsp, err := cs.transport.Request(req)
	if err != nil {
		return nil, err
	}

	result, err := rsp.ReadResult()
	if err != nil {
		return nil, err
	}
	var contents []Content
	err = unmarshalContents(result, &contents)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func (cs *clientSession) SubscribeResource(resource *ResourceDefinition) (<-chan *Notification, error) {
	req := &Request{
		Method: MethodSubscribeResource,
		Params: Params("{\"uri\": " + resource.URI + "}"),
	}
	_, err := cs.transport.Request(req)
	if err != nil {
		return nil, err
	}

	cs.subscribingResourcesLock.Lock()
	defer cs.subscribingResourcesLock.Unlock()

	ch := make(chan *Notification)

	cs.subscribingResources[resource.URI] = ch

	return ch, nil
}

func (s *clientSession) ListPrompts(ctx context.Context) ([]*PromptDefinition, error) {
	req := &Request{
		Method: MethodListPrompts,
	}
	rsp, err := s.transport.Request(req)
	if err != nil {
		return nil, err
	}
	result, err := rsp.ReadResult()
	if err != nil {
		return nil, err
	}
	var prompts []*PromptDefinition
	err = json.Unmarshal(result, &prompts)
	if err != nil {
		return nil, err
	}
	return prompts, nil
}

func (s *clientSession) GetPrompt(ctx context.Context, prompt *PromptDefinition, args map[string]any) ([]Content, error) {
	params := map[string]any{
		"name":      prompt.Name,
		"arguments": args,
	}
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	req := &Request{
		Method: MethodGetPrompt,
		Params: Params(paramsJSON),
	}
	rsp, err := s.transport.Request(req)
	if err != nil {
		return nil, err
	}
	result, err := rsp.ReadResult()
	if err != nil {
		return nil, err
	}
	var messages []struct {
		Role    string  `json:"role"`
		Content Content `json:"content"`
	}
	err = json.Unmarshal(result, &messages)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

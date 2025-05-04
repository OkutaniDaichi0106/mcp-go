package mcp

import (
	"context"
	"encoding/json"
)

type ServerSession interface {
	Close() error
	Shutdown() error
	// Notify() error

	//
	ListRoots(ctx context.Context) ([]*RootDefinition, error)

	Sample(ctx context.Context) error
}

var _ ServerSession = (*serverSession)(nil)

type serverSession struct {
	sessionID string

	transport Transport

	clientCapabilities Capabilities
}

func (s *serverSession) Close() error {
	return nil
}

func (s *serverSession) Shutdown() error {
	return nil
}

func (ss *serverSession) ListRoots(ctx context.Context) ([]*RootDefinition, error) {
	req := &Request{
		Method: MethodListRoots,
	}
	rsp, err := ss.transport.Request(req)
	if err != nil {
		return nil, err
	}

	result, err := rsp.ReadResult()
	if err != nil {
		return nil, err
	}

	var roots []*RootDefinition
	err = json.Unmarshal(result, &roots)
	if err != nil {
		return nil, err
	}

	return roots, nil
}

func (s *serverSession) Sample(ctx context.Context) error {
	return nil
}

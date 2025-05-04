package mcp

import (
	"context"
	"encoding/json"
)

type Transport interface {
	// Init(version Version, capabilities Capabilities, info map[string]any, options ...map[string]any) error

	Close() error
	CloseWithError(err error) error

	Request(req *Request) (ResponseReader, error)
	RequestSync(ctx context.Context, req *Request) (ResponseReader, error)

	Notify(notif *Notification) error

	AcceptRequest(ctx context.Context) (*Request, ResponseWriter, error)
	AcceptNotification(ctx context.Context) (*Notification, error)
}

type ResponseWriter interface {
	WriteResult(result Result) error
	CloseWithError(code ErrorCode, msg string, data map[string]json.RawMessage) error
}

type ResponseReader interface {
	ReadResult() (Result, error)
}

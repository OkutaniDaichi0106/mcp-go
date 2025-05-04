package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os/exec"

	"golang.org/x/exp/slog"
)

type Client struct {
	Name    string
	Version Version

	AdditionalClientInfo   map[string]any
	AdditionalCapabilities Capabilities

	RootsChangedNotification bool

	Handler ClientHandler

	sessionIDs map[string]int
	sessions   []*clientSession

	closed bool
	// closedErr error
}

func NewClient() *Client {
	c := &Client{}
	return c
}

func (c *Client) Dial(t Transport) (ClientSession, error) {
	return c.dial(context.Background(), t)
}

func (c *Client) DialWithContext(ctx context.Context, t Transport) (ClientSession, error) {
	return c.dial(ctx, t)
}

func (c *Client) DialStdio(command string, args ...string) (ClientSession, error) {
	cmd := exec.Command(command, args...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		slog.Error("failed to start command", "error", err)
		return nil, err
	}

	go func() {
		err := cmd.Wait()
		if err != nil {
			slog.Error("command exited with error", "error", err)
		}
	}()

	return c.DialStream(stdin, stdout)
}

func (c *Client) DialStream(w io.WriteCloser, r io.ReadCloser) (ClientSession, error) {
	return c.Dial(NewStreamTransport(w, r))
}

func (c *Client) DialHTTP(url string, config *HTTPConfig) (ClientSession, error) {
	transport := newClientHTTPTransport(url, config)

	// Set session ID
	key := SessionID("sessionID")
	ctx := context.WithValue(context.Background(), key, config.SessionID)

	return c.DialWithContext(ctx, transport)
}

func (c *Client) dial(ctx context.Context, t Transport) (*clientSession, error) {
	params := map[string]any{
		"protocolVersion": DefaultVersion,
		"capabilities":    c.capabilities(),
		"clientInfo":      c.info(),
	}

	paramsJSON, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}
	req := &Request{
		Method: MethodInit,
		Params: Params(paramsJSON),
	}

	r, err := t.RequestSync(ctx, req)
	if err != nil {
		return nil, err
	}

	result, err := r.ReadResult()
	if err != nil {
		return nil, err
	}
	var resultMapping map[string]json.RawMessage
	err = json.Unmarshal(result, &resultMapping)
	if err != nil {
		return nil, err
	}

	if Version(resultMapping["protocolVersion"]) != DefaultVersion {
		return nil, errors.New("protocol version mismatch")
	}

	var capabilities Capabilities
	err = json.Unmarshal(resultMapping["capabilities"], &capabilities)
	if err != nil {
		return nil, err
	}

	// TODO: Handle server serverInfo
	var serverInfo map[string]any
	err = json.Unmarshal(resultMapping["serverInfo"], &serverInfo)
	if err != nil {
		return nil, err
	}

	// Listen requests and handle them
	ctx, cancel := context.WithCancel(ctx)
	go c.handleRequests(ctx, t)

	sess := &clientSession{
		transport:            t,
		serverCapabilities:   capabilities,
		serverInfo:           serverInfo,
		subscribingResources: make(map[string]chan *Notification),
		cancelFunc:           cancel,
	}

	c.sessions = append(c.sessions, sess)

	if v := ctx.Value("sessionID"); v != nil {
		sessionID, ok := v.(string)
		if !ok {
			return nil, errors.New("sessionID is not string")
		}
		c.sessionIDs[sessionID] = len(c.sessions) - 1
	}

	return sess, nil
}

func (c *Client) handleRequests(ctx context.Context, t Transport) {
	// TODO: Implement
	for {
		req, w, err := t.AcceptRequest(ctx)
		if err != nil {
			return
		}

		switch req.Method {
		case MethodCreateSampleMessage:
			c.Handler.ServeSample(newContentsWriter(w), "", nil)
		case MethodNotifyRootChanged:
			c.Handler.ServeRootsChanged(t)
		}
	}
}

func (c *Client) Close() error {
	// Close all sessions
	for _, session := range c.sessions {
		session.Close()
	}

	c.closed = true

	return nil
}

func (c *Client) capabilities() Capabilities {
	capabilities := c.AdditionalCapabilities
	if capabilities == nil {
		capabilities = make(Capabilities)
	}

	capabilities.Merge(map[string]map[string]bool{
		"roots": {
			"listChanged": c.RootsChangedNotification,
		},
		"sampling": {},
	})

	return capabilities
}

func (c *Client) info() map[string]any {
	info := c.AdditionalClientInfo
	if info == nil {
		info = make(map[string]any)
	}

	info["name"] = c.Name
	info["version"] = c.Version

	return info
}

package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
)

type Server struct {
	Name                 string
	Version              string
	AdditionalServerInfo map[string]any

	// Capabilities related to resources
	ResourceSubscription         bool
	ResourcesChangedNotification bool
	// Capabilities related to prompts
	PromptsChangedNotification bool
	// Capabilities related to tools
	ToolsChangedNotification bool
	// Capabilities related to roots
	RootsChangedNotification bool
	// Additional capabilities
	AdditionalCapabilities Capabilities

	Options map[string]any

	Handler ServerHandler

	// http.Server
	Logger *slog.Logger

	initOnce    sync.Once
	initialized bool

	cancelFuncs     []context.CancelFunc
	cancelFuncsLock sync.Mutex

	// cancelFunc context.CancelFunc
	sessions map[string]*serverSession
}

func NewServer(name, version string) *Server {
	s := &Server{
		Name:    name,
		Version: version,
	}

	s.init()

	return s
}

func (s *Server) init() {
	s.initOnce.Do(func() {
		s.sessions = make(map[string]*serverSession)
		s.cancelFuncs = make([]context.CancelFunc, 0)

		s.initialized = true
	})
}

func (s *Server) AcceptStdio() (ServerSession, error) {
	if !s.initialized {
		s.init()
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	defer signal.Stop(signalCh)

	go func() {
		<-signalCh
		// TODO: close connection
		// s.cancelFunc
		s.Close()
	}()

	return s.AcceptStream(os.Stdout, os.Stdin)
}

func (s *Server) AcceptStream(w io.WriteCloser, r io.ReadCloser) (ServerSession, error) {
	if !s.initialized {
		s.init()
	}

	t := NewStreamTransport(w, r)

	return s.Accept(t)
}

func (s *Server) AcceptHTTP(w http.ResponseWriter, r *http.Request) (ServerSession, error) {
	if !s.initialized {
		s.init()
	}

	if r.Method == http.MethodPost {
		body := r.Body

		w.Header().Set("Content-Type", "application/json")
		sessionID := r.Header.Get("mcp-session-id")

		if sessionID == "" {
			// Handle just requests
			transport := newServerHTTPTransport()
			go s.handleRequests(context.Background(), transport)

			return nil, errors.New("session not found")
		}

		// Get the existing session
		session, ok := s.sessions[sessionID]
		if ok {
			// Get the existing transport
			transport, ok := session.transport.(*httpServerTransport)
			if !ok {
				return nil, errors.New("session is not HTTP session")
			}

			// Read the messages
			go transport.readMessages(w, body)

			return session, nil
		}

		// Create a new session
		transport := newServerHTTPTransport()
		go transport.readMessages(w, body)

		session, err := s.accept(transport)
		if err != nil {
			return nil, err
		}

		// Set the session ID
		session.sessionID = sessionID

		// Store the session
		s.sessions[sessionID] = session

		return session, nil
	} else if r.Method == http.MethodGet {
		sessionID := r.Header.Get("mcp-session-id")

		session, ok := s.sessions[sessionID]
		if !ok {
			return nil, errors.New("session not found")
		}

		httpTransport, ok := session.transport.(*httpServerTransport)
		if !ok {
			return nil, errors.New("session is not HTTP session")
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		if !ok {
			s.Logger.Error("failed to get flusher")
			return nil, errors.New("failed to get flusher")
		}

		// Set the SSE writer and flusher
		httpTransport.sseWriter = w
		httpTransport.sseFlusher = flusher

		return session, nil
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return nil, errors.New("method not allowed")
	}
}

func (s *Server) accept(t Transport) (*serverSession, error) {
	ctx := context.Background()
	ctx, cancelFunc := context.WithCancel(ctx)
	s.cancelFuncsLock.Lock()
	s.cancelFuncs = append(s.cancelFuncs, cancelFunc)
	s.cancelFuncsLock.Unlock()

	// Initialize
	req, w, err := t.AcceptRequest(ctx)
	if err != nil {
		return nil, err
	}

	if req.Method != MethodInit {
		return nil, errors.New("first request must be init")
	}

	var params map[string]json.RawMessage
	err = json.Unmarshal(req.Params, &params)
	if err != nil {
		return nil, err
	}

	if Version(params["protocolVersion"]) != DefaultVersion {
		return nil, errors.New("protocol version mismatch")
	}

	var capabilities Capabilities
	err = json.Unmarshal(params["capabilities"], &capabilities)
	if err != nil {
		return nil, err
	}

	// Write result
	result := s.Options
	result["protocolVersion"] = DefaultVersion
	result["capabilities"] = s.capabilities()
	result["serverInfo"] = s.info()

	resultJson, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	err = w.WriteResult(resultJson)
	if err != nil {
		return nil, err
	}

	session := &serverSession{
		transport:          t,
		clientCapabilities: capabilities,
	}

	// Listen requests and handle them
	go s.handleRequests(ctx, t)

	return session, nil
}

func (s *Server) Accept(t Transport) (ServerSession, error) {
	if !s.initialized {
		s.init()
	}

	return s.accept(t)
}

func (s *Server) Close() error {
	s.cancelFuncsLock.Lock()
	defer s.cancelFuncsLock.Unlock()

	for _, cancel := range s.cancelFuncs {
		cancel()
	}
	return nil
}

func (s *Server) capabilities() Capabilities {
	// Capabilities
	capabilities := s.AdditionalCapabilities
	if capabilities == nil {
		capabilities = make(Capabilities)
	}

	capabilities.Merge(map[string]map[string]bool{
		"tools": {
			"listChanged": s.ToolsChangedNotification,
		},
		"resources": {
			"listChanged": s.ResourcesChangedNotification,
			"subscribe":   s.ResourceSubscription,
		},
		"prompts": {
			"listChanged": s.PromptsChangedNotification,
		},
	})
	return capabilities
}

func (s *Server) info() map[string]any {
	info := s.AdditionalServerInfo
	if info == nil {
		info = make(map[string]any)
	}

	info["name"] = s.Name
	info["version"] = s.Version

	return info
}

func (s *Server) handleRequests(ctx context.Context, t Transport) {
	for {
		req, w, err := t.AcceptRequest(ctx)
		if err != nil {
			return
		}

		switch req.Method {
		case MethodListTools:
			// List tools
			listsJSON, err := json.Marshal(s.Handler.ListTools())
			if err != nil {
				slog.Error("failed to marshal tools", "error", err)
				continue
			}

			result := Result(`{"tools": ` + string(listsJSON) + `}`)

			// Write result
			err = w.WriteResult(result)
			if err != nil {
				continue
			}
		case MethodCallTool:
			// Call tool
			var params map[string]json.RawMessage
			err := json.Unmarshal(req.Params, &params)
			if err != nil {
				s.Logger.Error("failed to unmarshal params", "error", err)
				continue
			}
			name := string(params["name"])

			var args map[string]any
			err = json.Unmarshal(params["arguments"], &args)
			if err != nil {
				s.Logger.Error("failed to unmarshal params", "error", err)
				continue
			}

			s.Handler.ServeTool(newContentsWriter(w), name, args)
		case MethodListResources:
			// List resources
			result := map[string]any{
				"resources": s.Handler.ListResources(),
			}
			resultJson, err := json.Marshal(result)
			if err != nil {
				s.Logger.Error("failed to marshal resources", "error", err)
				continue
			}

			// Write result
			err = w.WriteResult(resultJson)
			if err != nil {
				s.Logger.Error("failed to write result", "error", err)
				continue
			}
		case MethodReadResource:
			// Read resource
			var params map[string]json.RawMessage
			err := json.Unmarshal(req.Params, &params)
			if err != nil {
				s.Logger.Error("failed to unmarshal params", "error", err)
				continue
			}

			if params["uri"] == nil {
				s.Logger.Error("missing uri field")
				continue
			}

			s.Handler.ServeResource(newContentsWriter(w), string(params["uri"]))

		case MethodSetLogLevel:
			// Set log level
			var params map[string]json.RawMessage
			err := json.Unmarshal(req.Params, &params)
			if err != nil {
				s.Logger.Error("failed to unmarshal params", "error", err)
				continue
			}

			// level := convertStrToLevel(params["level"].(string))
			// logger := slog.New(NewLogHandler(t, level))
		}
	}
}

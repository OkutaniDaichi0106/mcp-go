package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"golang.org/x/exp/slog"
)

func newClientHTTPTransport(url string, config *HTTPConfig) *httpClientTransport {
	if config == nil {
		config = &HTTPConfig{}
	}
	idGen := NewIDGenerator()
	t := &httpClientTransport{
		endpoint:                   url,
		client:                     &http.Client{Transport: config.Transport},
		generateIDFunc:             idGen.Generate,
		sentRequestMap:             make(map[ID]chan *response),
		rceivedRequestIDMap:        make(map[ID]struct{}),
		receivedRequestQueue:       newRequestsQueue(),
		receivedNotificationsQueue: newNotificationsQueue(),
		config:                     config,
	}

	go t.listenMessages()

	return t
}

var _ Transport = (*httpClientTransport)(nil)

type httpClientTransport struct {
	endpoint string
	client   *http.Client

	config *HTTPConfig

	sseReader io.ReadCloser

	generateIDFunc       func() ID
	sentRequestMap       map[ID]chan *response
	sentRequestIDMapLock sync.RWMutex

	rceivedRequestIDMap    map[ID]struct{}
	rceivedRequestsMapLock sync.RWMutex

	receivedRequestQueue       *requestsQueue
	receivedNotificationsQueue *notificationsQueue

	closed    bool
	closedErr error
}

// func (t *httpClientTransport) Init(params Params) (Result, error) {
// 	req := &Request{
// 		Method: MethodInit,
// 		Params: params,
// 	}
// 	res, err := t.RequestSync(context.TODO(), req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	result, err := res.ReadResult()
// 	if err != nil {
// 		return nil, err
// 	}

// 	sseReq, err := http.NewRequest("GET", t.endpoint, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	sseReq.Header.Set("Accept", "text/event-stream")
// 	if t.config.SessionID != "" {
// 		sseReq.Header.Set("mcp-session-id", string(t.config.SessionID))
// 	}

// 	sseResp, err := t.client.Do(sseReq)
// 	if err != nil {
// 		return nil, err
// 	}

// 	t.sseReader = sseResp.Body

// 	go t.listenMessages()

// 	return result, nil
// }

func (t *httpClientTransport) Close() error {
	if t.closed {
		if t.closedErr != nil {
			return fmt.Errorf("transport is already closed: %w", t.closedErr)
		}
		return errors.New("transport is already closed")
	}

	t.receivedRequestQueue.Clear()
	t.receivedNotificationsQueue.Clear()

	closeErr := t.sseReader.Close()
	if closeErr != nil {
		return closeErr
	}

	t.closed = true

	return nil
}

func (t *httpClientTransport) CloseWithError(err error) error {
	if t.closed {
		if t.closedErr != nil {
			return fmt.Errorf("transport is already closed: %w", t.closedErr)
		}
		return errors.New("transport is already closed")
	}

	t.receivedRequestQueue.Clear()
	t.receivedNotificationsQueue.Clear()

	closeErr := t.sseReader.Close()
	if closeErr != nil {
		return closeErr
	}

	t.closed = true
	t.closedErr = err

	return nil
}

func (t *httpClientTransport) Request(req *Request) (ResponseReader, error) {
	id := t.generateIDFunc()
	if _, ok := t.sentRequestMap[id]; ok {
		panic("ID is already used")
	}

	rspCh := make(chan *response)

	// Record request ID
	t.sentRequestIDMapLock.Lock()
	t.sentRequestMap[id] = rspCh
	t.sentRequestIDMapLock.Unlock()

	jsonrpcReq := map[string]any{
		"jsonrpc": "2.0",
		"method":  req.Method,
		"params":  req.Params,
		"id":      id,
	}
	body, err := json.Marshal(jsonrpcReq)
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequest("POST", t.endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json, text/event-stream")
	if t.config.SessionID != "" {
		httpReq.Header.Set("mcp-session-id", string(t.config.SessionID))
	}
	// httpReq.Header.Set("Authorization", "Bearer "+req.token)

	go func() {
		defer close(rspCh)

		httpResp, err := t.client.Do(httpReq)
		if err != nil {
			rspCh <- &response{
				result: nil,
				err: &Error{
					Code:    JSONRPCInternalErrorCode,
					Message: err.Error(),
				},
			}
			return
		}

		if httpResp.StatusCode != http.StatusOK {
			rspCh <- &response{
				result: nil,
				err: &Error{
					Code:    JSONRPCInternalErrorCode,
					Message: "HTTP status code is not 200",
					Data: map[string]any{
						"status_code": httpResp.StatusCode,
					},
				},
			}
			return
		}

		defer httpResp.Body.Close()

		var jsonrpcResp map[string]json.RawMessage
		err = json.NewDecoder(httpResp.Body).Decode(&jsonrpcResp)
		if err != nil {
			rspCh <- &response{
				result: nil,
				err: &Error{
					Code:    JSONRPCInternalErrorCode,
					Message: err.Error(),
				},
			}
			return
		}

		if id != ID(jsonrpcResp["id"]) {
			rspCh <- &response{
				result: nil,
				err: &Error{
					Code:    JSONRPCInternalErrorCode,
					Message: "unexpected request ID",
				},
			}
			return
		}

		if jsonrpcResp["error"] != nil {
			var errObj Error
			err = json.Unmarshal(jsonrpcResp["error"], &errObj)
			if err != nil {
				rspCh <- &response{
					result: nil,
					err: &Error{
						Code:    JSONRPCInternalErrorCode,
						Message: fmt.Sprintf("failed to unmarshal error: %s", err.Error()),
						Data: map[string]any{
							"error": jsonrpcResp["error"],
						},
					},
				}
				return
			}

			rspCh <- &response{
				result: nil,
				err:    &errObj,
			}
		}

		if jsonrpcResp["result"] == nil {
			rspCh <- &response{
				result: nil,
				err: &Error{
					Code:    JSONRPCInternalErrorCode,
					Message: "missing result field",
				},
			}
			return
		}

		rspCh <- &response{
			result: Result(jsonrpcResp["result"]),
			err:    nil,
		}
	}()

	return &asyncResponseReader{ch: rspCh}, nil
}

func (t *httpClientTransport) RequestSync(ctx context.Context, req *Request) (ResponseReader, error) {
	id := t.generateIDFunc()
	if _, ok := t.sentRequestMap[id]; ok {
		panic("ID is already used")
	}

	// Record request ID
	t.sentRequestIDMapLock.Lock()
	t.sentRequestMap[id] = nil
	t.sentRequestIDMapLock.Unlock()

	jsonrpcReq := map[string]any{
		"jsonrpc": "2.0",
		"method":  req.Method,
		"params":  req.Params,
		"id":      id,
	}
	body, err := json.Marshal(jsonrpcReq)
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, "POST", t.endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json, text/event-stream")
	// httpReq.Header.Set("Authorization", "Bearer "+req.token)

	httpResp, err := t.client.Do(httpReq)
	if err != nil {
		return nil, err
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status code is not 200: %d", httpResp.StatusCode)
	}

	defer httpResp.Body.Close()

	var jsonrpcResp map[string]json.RawMessage
	err = json.NewDecoder(httpResp.Body).Decode(&jsonrpcResp)
	if err != nil {
		return nil, err
	}

	if id != ID(jsonrpcResp["id"]) {
		return nil, errors.New("unexpected request ID")
	}

	if jsonrpcResp["error"] != nil {
		var errObj Error
		err = json.Unmarshal(jsonrpcResp["error"], &errObj)
		if err != nil {
			return nil, err
		}

		return &response{
			result: nil,
			err:    &errObj,
		}, nil
	}

	if jsonrpcResp["result"] == nil {
		return nil, errors.New("missing result field")
	}

	return &response{
		result: Result(jsonrpcResp["result"]),
		err:    nil,
	}, nil
}

func (t *httpClientTransport) Notify(notif *Notification) error {
	jsonrpcReq := map[string]any{
		"jsonrpc": "2.0",
		"method":  notif.Method,
		"params":  notif.Params,
	}
	body, err := json.Marshal(jsonrpcReq)
	if err != nil {
		return err
	}
	httpReq, err := http.NewRequest("POST", t.endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json, text/event-stream")
	// httpReq.Header.Set("Authorization", "Bearer "+req.token)

	httpResp, err := t.client.Do(httpReq)
	if err != nil {
		return err
	}

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP status code is not 200: %d", httpResp.StatusCode)
	}

	defer httpResp.Body.Close()

	// Ignore response

	return nil
}

func (t *httpClientTransport) AcceptRequest(ctx context.Context) (*Request, ResponseWriter, error) {
	for {
		if t.receivedRequestQueue.Len() > 0 {
			return t.receivedRequestQueue.Pop(), nil, nil
		}

		select {
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		case <-t.receivedRequestQueue.Chan():
			continue
		}
	}
}

func (t *httpClientTransport) AcceptNotification(ctx context.Context) (*Notification, error) {
	for {
		if t.receivedNotificationsQueue.Len() > 0 {
			return t.receivedNotificationsQueue.Pop(), nil
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-t.receivedNotificationsQueue.Chan():
			continue
		}
	}
}

func (t *httpClientTransport) listenMessages() {
	decoder := json.NewDecoder(t.sseReader)
	for {

		// Read a message as single JSON object
		var message map[string]json.RawMessage
		err := decoder.Decode(&message)
		if err == nil {
			t.listenSingleMessage(message)
			continue
		}

		// Read a message as batch JSON object
		var batch []map[string]json.RawMessage
		err = decoder.Decode(&batch)
		if err == nil {
			for _, message := range batch {
				t.listenSingleMessage(message)
			}
			continue
		}

		err = errors.New("the message is neither single nor batch")
		slog.Error("failed to decode message", "error", err)
	}

}

func (t *httpClientTransport) listenSingleMessage(message map[string]json.RawMessage) {
	id, hasID := message["id"]
	rawMethod, hasMethod := message["method"]
	if hasID && hasMethod {
		// It should be a request

		// Check if params field is present
		rawParams, hasParams := message["params"]
		if !hasParams {
			slog.Error("Missing params field for request")
			return
		}

		// Check if ID is valid
		if string(id) == "" {
			slog.Error("Invalid request ID", "id", id)
			return
		}

		// Check if ID is already received
		t.rceivedRequestsMapLock.RLock()
		if _, ok := t.rceivedRequestIDMap[ID(id)]; ok {
			slog.Error("Duplicate request ID", "id", id)
			return
		}
		t.rceivedRequestsMapLock.RUnlock()

		// Record received request ID
		t.rceivedRequestsMapLock.Lock()
		t.rceivedRequestIDMap[ID(id)] = struct{}{}
		t.rceivedRequestsMapLock.Unlock()

		// Push request to queue
		t.receivedRequestQueue.Push(&Request{
			Method: Method(rawMethod),
			Params: Params(rawParams),
		})
	} else if !hasID && hasMethod {
		// It should be a notification
		// Check if params field is present
		rawParams, hasParams := message["params"]
		if !hasParams {
			slog.Error("Missing params field for request")
			return
		}

		// Check if ID is valid
		if string(id) == "" {
			slog.Error("Invalid request ID", "id", id)
			return
		}

		// Queue notification
		t.receivedNotificationsQueue.Push(&Notification{
			Method: Method(rawMethod),
			Params: Params(rawParams),
		})
	} else if hasID && !hasMethod {
		// It should be a response

		t.sentRequestIDMapLock.RLock()
		rspCh, ok := t.sentRequestMap[ID(id)]
		if !ok {
			slog.Error("Unknown request ID", "id", id)
			return
		}
		t.sentRequestIDMapLock.RUnlock()

		rawResult, hasResult := message["result"]
		rawError, hasError := message["error"]
		if hasError {
			errObj := Error{}
			err := json.Unmarshal(rawError, &errObj)
			if err != nil {
				slog.Error("Failed to unmarshal error", "error", err)
				return
			}

			rspCh <- &response{
				result: nil,
				err:    &errObj,
			}

			close(rspCh)
		} else if hasResult {
			rspCh <- &response{
				result: Result(rawResult),
				err:    nil,
			}

			close(rspCh)
		}

	} else {
		slog.Error("Unknown message type")
		return
	}
}

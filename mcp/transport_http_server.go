package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"

	"golang.org/x/exp/slog"
)

var _ Transport = (*httpServerTransport)(nil)

func newServerHTTPTransport() *httpServerTransport {
	idGen := NewIDGenerator()
	t := &httpServerTransport{
		generateIDFunc:             idGen.Generate,
		sentRequestMap:             make(map[ID]chan *response),
		rceivedRequestIDMap:        make(map[ID]http.ResponseWriter),
		receivedRequestQueue:       newRequestsQueue(),
		receivedNotificationsQueue: newNotificationsQueue(),
	}

	return t
}

type httpServerTransport struct {
	generateIDFunc func() ID

	sentRequestMap       map[ID]chan *response
	sentRequestIDMapLock sync.RWMutex

	sseWriter  http.ResponseWriter
	sseFlusher http.Flusher
	wmu        sync.Mutex

	rceivedRequestIDMap    map[ID]http.ResponseWriter
	rceivedRequestsMapLock sync.RWMutex

	receivedRequestQueue       *requestsQueue
	receivedNotificationsQueue *notificationsQueue

	closed    bool
	closedErr error
}

func (t *httpServerTransport) Close() error {
	t.closed = true
	return nil
}

func (t *httpServerTransport) CloseWithError(err error) error {
	t.closed = true

	t.wmu.Lock()
	defer t.wmu.Unlock()

	t.sseWriter.WriteHeader(http.StatusInternalServerError)
	t.sseWriter.Write([]byte("event: error\ndata: " + err.Error() + "\n\n"))
	t.sseFlusher.Flush()

	t.closedErr = err

	return nil
}

func (t *httpServerTransport) Request(req *Request) (ResponseReader, error) {
	id := t.generateIDFunc()
	if _, ok := t.sentRequestMap[id]; ok {
		panic("ID is already used")
	}

	rspCh := make(chan *response)

	t.sentRequestIDMapLock.Lock()
	t.sentRequestMap[id] = rspCh
	t.sentRequestIDMapLock.Unlock()

	jsonrpcReq := map[string]any{
		"jsonrpc": "2.0",
		"method":  req.Method,
		"params":  req.Params,
		"id":      id,
	}

	t.wmu.Lock()

	err := json.NewEncoder(t.sseWriter).Encode(jsonrpcReq)
	t.sseFlusher.Flush()
	t.wmu.Unlock()
	if err != nil {
		return nil, err
	}

	return &asyncResponseReader{ch: rspCh}, nil
}

func (t *httpServerTransport) RequestSync(ctx context.Context, req *Request) (ResponseReader, error) {
	id := t.generateIDFunc()
	if _, ok := t.sentRequestMap[id]; ok {
		panic("ID is already used")
	}

	rspCh := make(chan *response)

	t.sentRequestIDMapLock.Lock()
	t.sentRequestMap[id] = rspCh
	t.sentRequestIDMapLock.Unlock()

	jsonrpcReq := map[string]any{
		"jsonrpc": "2.0",
		"method":  req.Method,
		"params":  req.Params,
		"id":      id,
	}

	t.wmu.Lock()
	err := json.NewEncoder(t.sseWriter).Encode(jsonrpcReq)
	t.sseFlusher.Flush()
	t.wmu.Unlock()

	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case rsp := <-rspCh:
		return rsp, nil
	}
}

func (t *httpServerTransport) Notify(notif *Notification) error {
	return nil
}

func (t *httpServerTransport) AcceptRequest(ctx context.Context) (*Request, ResponseWriter, error) {
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

func (t *httpServerTransport) AcceptNotification(ctx context.Context) (*Notification, error) {
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

func (t *httpServerTransport) readMessages(w http.ResponseWriter, r io.ReadCloser) {
	decoder := json.NewDecoder(r)
	for {
		// Read a message as single JSON object
		var message map[string]json.RawMessage
		err := decoder.Decode(&message)
		if err == nil {
			t.readSingleMessage(w, message)
			continue
		}

		// Read a message as batch JSON object
		var batch []map[string]json.RawMessage
		err = decoder.Decode(&batch)
		if err == nil {
			for _, message := range batch {
				t.readSingleMessage(w, message)
			}
			continue
		}

		err = errors.New("the message is neither single nor batch")
		slog.Error("failed to decode message", "error", err)
	}
}

func (t *httpServerTransport) readSingleMessage(w http.ResponseWriter, message map[string]json.RawMessage) {
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
		t.rceivedRequestIDMap[ID(id)] = w
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

package mcp

type requestsQueue struct {
	queue []*Request

	ch chan struct{}
}

func newRequestsQueue() *requestsQueue {
	return &requestsQueue{
		queue: make([]*Request, 0),
		ch:    make(chan struct{}),
	}
}

func (q *requestsQueue) Push(req *Request) {
	q.queue = append(q.queue, req)
	q.ch <- struct{}{}
}

func (q *requestsQueue) Pop() *Request {
	req := q.queue[0]
	q.queue = q.queue[1:]
	return req
}

func (q *requestsQueue) Empty() bool {
	return len(q.queue) == 0
}

func (q *requestsQueue) Chan() chan struct{} {
	return q.ch
}

func (q *requestsQueue) Len() int {
	return len(q.queue)
}

func (q *requestsQueue) Clear() {
	q.queue = q.queue[:0]
}

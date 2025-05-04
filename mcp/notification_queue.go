package mcp

type notificationsQueue struct {
	queue []*Notification

	ch chan struct{}
}

func newNotificationsQueue() *notificationsQueue {
	return &notificationsQueue{
		queue: make([]*Notification, 0),
		ch:    make(chan struct{}),
	}
}

func (q *notificationsQueue) Push(req *Notification) {
	q.queue = append(q.queue, req)
	q.ch <- struct{}{}
}

func (q *notificationsQueue) Pop() *Notification {
	req := q.queue[0]
	q.queue = q.queue[1:]
	return req
}

func (q *notificationsQueue) Empty() bool {
	return len(q.queue) == 0
}

func (q *notificationsQueue) Chan() chan struct{} {
	return q.ch
}

func (q *notificationsQueue) Len() int {
	return len(q.queue)
}

func (q *notificationsQueue) Clear() {
	q.queue = q.queue[:0]
}

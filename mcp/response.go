package mcp

// type Responce interface {
// 	ID() string
// 	Result() any
// 	Encode(w io.Writer) error
// 	Decode(r io.Reader) error
// }

type response struct {
	result Result
	err    *Error
}

func (r *response) ReadResult() (Result, error) {
	if r.err != nil {
		return nil, r.err
	}

	return r.result, nil
}

type asyncResponseReader struct {
	ch chan *response
}

func (r *asyncResponseReader) ReadResult() (Result, error) {
	rsp := <-r.ch
	close(r.ch)

	if rsp.err != nil {
		return nil, rsp.err
	}

	return rsp.result, nil
}

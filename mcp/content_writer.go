package mcp

import (
	"errors"
	"fmt"
)

type ContentsWriter interface {
	WriteContents(content []Content) error
	CloseWithError(code ErrorCode, msg string) error
}

func newContentsWriter(rw ResponseWriter) ContentsWriter {
	return &contentsWriter{rw: rw}
}

var _ ContentsWriter = (*contentsWriter)(nil)

type contentsWriter struct {
	done      bool
	closedErr error

	rw ResponseWriter
}

func (cw *contentsWriter) WriteContents(contents []Content) error {
	if cw.done {
		if cw.closedErr != nil {
			return fmt.Errorf("writer is already closed: %w", cw.closedErr)
		}

		return errors.New("session has already done")
	}

	result, err := marshalContents(&contents)
	if err != nil {
		return err
	}

	err = cw.rw.WriteResult(Result(result))
	if err != nil {
		return err
	}

	cw.done = true

	return nil
}

func (cw *contentsWriter) CloseWithError(code ErrorCode, msg string) error {
	if cw.done {
		if cw.closedErr != nil {
			return fmt.Errorf("writer is already closed: %w", cw.closedErr)
		}

		return errors.New("session has already done")
	}

	return cw.rw.CloseWithError(code, msg, nil)
}

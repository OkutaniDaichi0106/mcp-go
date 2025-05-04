package mcp

import (
	"strconv"
	"sync/atomic"
)

type ID string

func NewIDGenerator() IDGenerator {
	return IDGenerator{}
}

type IDGenerator struct {
	counter atomic.Uint64
}

func (g *IDGenerator) Generate() ID {
	return ID(strconv.FormatUint(g.counter.Add(1), 10))
}

package iterator

import (
	"bufio"
	"bytes"
	"io"
)

type ReadLineIterator struct {
	bb  bytes.Buffer
	br  *bufio.Reader
	err error
}

func NewReadLineIterator(r io.Reader) *ReadLineIterator {
	return &ReadLineIterator{br: bufio.NewReader(r)}
}

func (r *ReadLineIterator) Next() bool {
	r.bb.Reset()
	for {
		line, prefix, err := r.br.ReadLine()
		if err != nil {
			r.err = err
			return false
		}
		r.bb.Write(line)
		if !prefix {
			break
		}
	}
	return true
}

func (r *ReadLineIterator) Err() error {
	return r.err
}

func (r *ReadLineIterator) Value() interface{} {
	return r.bb.Bytes()
}

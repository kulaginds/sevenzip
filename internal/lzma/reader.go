// Package lzma implements the LZMA decompressor.
package lzma

import (
	"bufio"
	"errors"
	"fmt"
	"io"

	"github.com/kulaginds/lzma"
)

type readCloser struct {
	c io.Closer
	r io.Reader
}

var (
	errAlreadyClosed = errors.New("lzma: already closed")
	errNeedOneReader = errors.New("lzma: need exactly one reader")
)

func (rc *readCloser) Close() error {
	if rc.c == nil || rc.r == nil {
		return errAlreadyClosed
	}

	if err := rc.c.Close(); err != nil {
		return fmt.Errorf("lzma: error closing: %w", err)
	}

	rc.c, rc.r = nil, nil

	return nil
}

func (rc *readCloser) Read(p []byte) (int, error) {
	if rc.r == nil {
		return 0, errAlreadyClosed
	}

	n, err := rc.r.Read(p)
	if err != nil && !errors.Is(err, io.EOF) {
		err = fmt.Errorf("lzma: error reading: %w", err)
	}

	return n, err
}

// NewReader returns a new LZMA io.ReadCloser.
func NewReader(p []byte, s uint64, readers []io.ReadCloser) (io.ReadCloser, error) {
	if len(readers) != 1 {
		return nil, errNeedOneReader
	}

	lr, err := lzma.NewReader1ForSevenZip(bufio.NewReader(readers[0]), p, s)
	if err != nil {
		return nil, fmt.Errorf("lzma: error creating reader: %w", err)
	}

	return &readCloser{
		c: readers[0],
		r: lr,
	}, nil
}

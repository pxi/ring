// Package ring implements an in-memory circular buffers.
package ring

import "io"

// NewWriter returns a new io.Writer that writes to p starting from offset off.
// It will panic if len(p) == 0.
func NewWriter(p []byte, off int64) io.Writer { return newRing(p, off, -1) }

// NewLimitWriter returns a new io.Writer that writes at most n bytes to p
// starting from offset off. It will panic if len(p) == 0.
func NewLimitWriter(p []byte, off, n int64) io.Writer { return newRing(p, off, n) }

// NewReader returns a new io.Reader that reads from p starting from offset off.
// It will panic if len(p) == 0.
func NewReader(p []byte, off int64) io.Reader { return newRing(p, off, -1) }

// NewLimitReader returns a new io.Reader that reads from p at most n bytes
// starting from offset off. It will panic if len(p) == 0.
func NewLimitReader(p []byte, off, n int64) io.Reader { return newRing(p, off, n) }

type ring struct {
	buf []byte
	off int64
	n   int64 // bytes remaining
}

func newRing(p []byte, off, n int64) *ring {
	if len(p) == 0 {
		panic("ring: empty buffer")
	}
	return &ring{p, off, n}
}

// Write writes up to len(p) bytes from p to underlying buffer.
func (r *ring) Write(p []byte) (n int, err error) {
	i := r.off
	s := int64(len(r.buf))
	if r.n > -1 && s > r.n {
		err = io.ErrShortWrite
		p = p[:r.n]
	}
	for n < len(p) {
		m := copy(r.buf[i:], p[n:])
		n += m
		i += int64(m)
		if i == s {
			i = 0
		}
	}
	r.off = i
	if r.n > -1 {
		r.n -= int64(n)
	}
	return
}

// Read reads up to len(p) bytes from underlying buffer to p.
func (r *ring) Read(p []byte) (n int, err error) {
	i := r.off
	s := int64(len(r.buf))
	if r.n > -1 && s > r.n {
		err = io.EOF
		p = p[:r.n]
	}
	for n < len(p) {
		m := copy(p[n:], r.buf[i:])
		n += m
		i += int64(m)
		if i == s {
			i = 0
		}
	}
	r.off = i
	if r.n > -1 {
		r.n -= int64(n)
	}
	return
}

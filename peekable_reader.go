package io

import "io"

type Peaker interface {
	io.Reader
	Peak(p []byte) (n int, err error)
}

type peaker struct {
	r    io.Reader
	peak []byte
}

func (p *peaker) Peak(b []byte) (n int, err error) {
	if len(p.peak) >= len(b) {
		n = copy(b, p.peak)
		return n, nil
	}

	n = copy(b, p.peak)
	nRead, err := p.r.Read(b[n:])
	n += nRead
	p.peak = make([]byte, n)
	copy(p.peak, b)
	return n, err
}

func (p *peaker) Read(b []byte) (n int, err error) {
	if len(p.peak) > 0 {
		n = copy(b, p.peak)
		p.peak = p.peak[n:]
		return n, err
	}

	return p.r.Read(b)
}

var _ Peaker = &peaker{}

//	p.peak = p.peak[n:]
func NewPeakableReader(r io.Reader) *peaker {
	return &peaker{r: r}
}

package parser

import (
	"io"
)

type eofReader struct {
	r     io.Reader
	state int
}

func (r *eofReader) Read(buf []byte) (n int, err error) {
	switch r.state {
	case 0:
		n, err = r.r.Read(buf)
		if err == io.EOF {
			r.state++
			err = nil
		}
		return

	case 1:
		r.state++
		buf[0] = '\n'
		return 1, nil

	default:
		return r.r.Read(buf)
	}
}

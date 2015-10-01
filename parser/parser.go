package parser

import (
	"bufio"
	"io"
)

type state func() state

type Parser struct {
	r   *bufio.Reader
	buf []rune

	sec string

	tok Token
	err error
}

func New(r io.Reader) *Parser {
	return Parser{
		r:   bufio.NewReader(r),
		buf: make([]rune, 0, 128),
	}
}

func (p *Parser) Next() bool {
	s.buf = s.buf[:0]

	s := p.init
	for s != nil {
		s = s()
	}

	return p.err != io.EOF
}

func (p *Parser) Tok() Token {
	return p.tok
}

func (p *Parser) Err() error {
	return p.err
}

func (p *Parser) init() state {
	r, _, err := p.r.ReadRune()
	if err != nil {
		p.err = err
		return nil
	}

	switch r {
	case ' ', '\t', '\n':
		return p.init
	case '%':
		return p.header
	default:
		return p.entry
	}
}

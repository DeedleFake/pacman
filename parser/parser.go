package parser

import (
	"bufio"
	"io"
	"unicode"
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
	return &Parser{
		r:   bufio.NewReader(&eofReader{r: r}),
		buf: make([]rune, 0, 128),
	}
}

func (p *Parser) Err() error {
	if p.err == io.EOF {
		return nil
	}

	return p.err
}

func (p *Parser) Next() bool {
	if p.err != nil {
		return false
	}

	p.buf = p.buf[:0]

	s := p.whitespace
	for s != nil {
		s = s()
	}

	return p.err == nil
}

func (p *Parser) Tok() Token {
	return p.tok
}

func (p *Parser) whitespace() state {
	r, _, err := p.r.ReadRune()
	if err != nil {
		p.err = err
		return nil
	}

	switch {
	case unicode.IsSpace(r):
		return p.whitespace
	case r == '%':
		return p.header
	default:
		p.r.UnreadRune()
		return p.entry
	}
}

func (p *Parser) header() state {
	r, _, err := p.r.ReadRune()
	if err != nil {
		p.err = err
		return nil
	}

	if r == '%' {
		p.tok = Header(p.buf)
		return nil
	}

	p.buf = append(p.buf, r)

	return p.header
}

func (p *Parser) entry() state {
	r, _, err := p.r.ReadRune()
	if err != nil {
		p.err = err
		return nil
	}

	if r == '\n' {
		p.tok = Entry(p.buf)
		return nil
	}

	p.buf = append(p.buf, r)

	return p.entry
}

type Token interface{}

type Header string

type Entry string

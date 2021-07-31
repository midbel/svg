package dot

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"unicode/utf8"
)

const (
	kwGraph   = "graph"
	kwInclude = "include"
)

type Graph struct {
}

type parser struct {
	scan *Scanner
	curr Token
	peek Token

	env map[string][]string
}

func Parse(r io.Reader) error {
	s, err := Scan(r)
	if err != nil {
		return err
	}
	p := parser{
		scan: s,
		env:  make(map[string][]string),
	}
	p.next()
	p.next()
	return p.parse()
}

func (p *parser) parse() error {
	var err error
	for err == nil && !p.done() {
		if p.curr.Type == comment {
			p.next()
			continue
		}
		if p.curr.Type != ident && p.curr.Type != lcurly {
			err = unexpectedToken(p.curr)
			break
		}
		switch {
		case p.peek.Type == assign:
			err = p.parseAssignment()
		case p.curr.Literal == kwInclude:
			err = p.parseInclude()
		case p.curr.Literal == kwGraph || p.curr.Type == lcurly:
			err = p.parseGraph()
		default:
			err = unexpectedToken(p.curr)
		}
		if err != nil {
			break
		}
	}
	return err
}

func (p *parser) parseAssignment() error {
	fmt.Println("enter parseAssignment")
	defer fmt.Println("leave parseAssignment")

	name := p.curr.Literal
	if _, ok := p.env[name]; ok {
		return fmt.Errorf("%s: already defined", name)
	}
	p.next()
	p.next()
	for p.curr.IsValue() {
		if p.curr.Type == variable {
			vs, ok := p.env[p.curr.Literal]
			if !ok {
				return fmt.Errorf("%s: undefined", p.curr.Literal)
			}
			p.env[name] = append(p.env[name], vs...)
		} else {
			p.env[name] = append(p.env[name], p.curr.Literal)
		}
		p.next()
	}
	fmt.Printf(">> assignment: %s = %s\n", name, p.env[name])
	return p.parseEOL()
}

func (p *parser) parseInclude() error {
	fmt.Println("enter parseInclude")
	defer fmt.Println("leave parseInclude")
	p.next()
	if p.curr.Type != text && p.curr.Type != ident {
		return unexpectedToken(p.curr)
	}
	fmt.Printf(">> include: %s\n", p.curr.Literal)
	p.next()
	return p.parseEOL()
}

func (p *parser) parseGraph() error {
	fmt.Println("enter parseGraph")
	defer fmt.Println("leave parseGraph")
	if p.curr.Type == ident {
		p.next()
		if p.curr.IsValue() {
			fmt.Printf(">> name: %s\n", p.curr.Literal)
			p.next()
		}
	}
	if p.curr.Type != lcurly {
		return unexpectedToken(p.curr)
	}
	p.next()
	for !p.done() && p.curr.Type != rcurly {
		if p.curr.Type == comment {
			p.next()
			continue
		}
		if p.curr.Type != ident {
			return unexpectedToken(p.curr)
		}
		var err error
		if p.peek.Type == equal {
			err = p.parseAttributes()
		} else {
			err = p.parseEdge()
		}
		if err != nil {
			return err
		}
	}
	if p.curr.Type != rcurly {
		return unexpectedToken(p.curr)
	}
	p.next()
	return p.parseEOL()
}

func (p *parser) parseEdge() error {
	fmt.Println("enter parseEdge")
	defer fmt.Println("leave parseEdge")
	for !p.curr.IsEOL() && p.curr.Type != comment && !p.done() {
		var nodes []string
		switch p.curr.Type {
		case lcurly:
			p.next()
			for p.curr.Type != rcurly {
				nodes = append(nodes, p.curr.Literal)
				p.next()
				if p.curr.Type == lsquare {
					p.next()
					if err := p.parseProperties(); err != nil {
						return err
					}
				}
			}
			if p.curr.Type != rcurly {
				return unexpectedToken(p.curr)
			}
			p.next()
			if p.curr.Type == lsquare {
				p.next()
				err := p.parseProperties()
				if err != nil {
					return err
				}
			}
		case ident:
			nodes = append(nodes, p.curr.Literal)
			p.next()
			if p.curr.Type == lsquare {
				p.next()
				err := p.parseProperties()
				if err != nil {
					return err
				}
				if p.curr.Type == comment || p.curr.IsEOL() {
					p.next()
					return nil
				}
			}
		default:
			return unexpectedToken(p.curr)
		}
		fmt.Printf(">> nodes: %s\n", nodes)
		if p.curr.Type == edge {
			if p.peek.Type != ident && p.peek.Type != lcurly {
				return unexpectedToken(p.peek)
			}
			p.next()
		}
	}
	return p.parseEOL()
}

func (p *parser) parseAttributes() error {
	fmt.Println("enter parseAttributes")
	defer fmt.Println("leave enter parseAttributes")
	var (
		name   = p.curr.Literal
		values []string
	)
	if p.peek.Type != equal {
		return unexpectedToken(p.curr)
	}
	p.next()
	p.next()
	for p.curr.IsValue() {
		values = append(values, p.curr.Literal)
		p.next()
	}
	if !p.curr.IsEOL() && p.curr.Type != comment {
		return unexpectedToken(p.curr)
	}
	fmt.Printf(">> attributes: %s = %s\n", name, values)
	return p.parseEOL()
}

func (p *parser) parseProperties() error {
	fmt.Println("enter parseProperties")
	defer fmt.Println("leave parseProperties")
	for !p.done() && p.curr.Type != rsquare {
		if p.curr.Type != ident && p.peek.Type != equal {
			return unexpectedToken(p.curr)
		}
		var (
			name   = p.curr.Literal
			values []string
		)
		p.next()
		p.next()
		for p.curr.IsValue() {
			values = append(values, p.curr.Literal)
			p.next()
		}
		if p.curr.Type != comma && p.curr.Type != rsquare {
			return unexpectedToken(p.curr)
		}
		if p.curr.Type == comma {
			p.next()
		}
		if p.curr.Type == comment {
			p.next()
		}
		fmt.Printf(">> properties: %s = %s\n", name, values)
	}
	if p.curr.Type != rsquare {
		return unexpectedToken(p.curr)
	}
	p.next()
	return nil
}

func (p *parser) parseEOL() error {
	if p.curr.Type == comment {
		p.next()
		return nil
	}
	if !p.curr.IsEOL() {
		return unexpectedToken(p.curr)
	}
	p.next()
	return nil
}

func (p *parser) next() {
	p.curr = p.peek
	p.peek = p.scan.Scan()
}

func (p *parser) done() bool {
	return p.curr.Type == eof
}

const (
	ident rune = -(iota + 1)
	text
	number
	variable
	assign
	edge
	comment
	invalid
)

type Token struct {
	Literal string
	Type    rune
}

func (t Token) IsEOL() bool {
	return t.Type == semicolon
}

func (t Token) IsEOF() bool {
	return t.Type == eof
}

func (t Token) IsInvalid() bool {
	return t.Type == invalid
}

func (t Token) IsValue() bool {
	return t.Type == number || t.Type == text || t.Type == ident || t.Type == variable
}

func (t Token) String() string {
	var prefix string
	switch t.Type {
	case ident:
		prefix = "ident"
	case text:
		prefix = "text"
	case number:
		prefix = "number"
	case variable:
		prefix = "variable"
	case comment:
		prefix = "comment"
	case assign:
		prefix = "assign"
	case edge:
		prefix = "edge"
	case invalid:
		prefix = "invalid"
	case eof:
		return "<eof>"
	case lcurly:
		return "<beg-obj>"
	case rcurly:
		return "<end-obj>"
	case lsquare:
		return "<beg-list>"
	case rsquare:
		return "<end-list>"
	case comma:
		return "<comma>"
	case semicolon:
		return "<semicolon>"
	case equal:
		return "<equal>"
	}
	return fmt.Sprintf("<%s(%s)>", prefix, t.Literal)
}

var ErrSyntax = errors.New("syntax error")

func unexpectedToken(tok Token) error {
	return fmt.Errorf("%w: unexpected token %s", ErrSyntax, tok)
}

type Scanner struct {
	buffer []byte
	curr   int
	next   int
	char   rune
}

func Scan(r io.Reader) (*Scanner, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var s Scanner
	s.buffer = bytes.ReplaceAll(buf, []byte{cr, nl}, []byte{nl})

	s.read()
	s.skipNL()
	return &s, nil
}

func (s *Scanner) Scan() Token {
	var t Token
	if s.char == eof {
		return t
	}
	if s.char == slash && s.peek() == slash {
		s.scanComment(&t)
		return t
	}
	if s.char == dash && s.peek() == dash {
		s.read()
		s.read()
		t.Type = edge
		return t
	}
	s.skipBlanks()
	switch {
	case isLetter(s.char):
		s.scanIdent(&t)
	case isDigit(s.char):
		s.scanNumber(&t)
	case isPunct(s.char):
		s.scanPunct(&t)
	case isQuote(s.char):
		s.scanText(&t)
	case isVariable(s.char):
		s.scanVariable(&t)
	case isNL(s.char):
		t.Type = semicolon
		s.skipNL()
	default:
		t.Type = invalid
	}
	s.skipBlanks()
	return t
}

func (s *Scanner) scanIdent(tok *Token) {
	pos := s.curr
	for isIdent(s.char) {
		s.read()
	}
	tok.Literal = string(s.buffer[pos:s.curr])
	tok.Type = ident
}

func (s *Scanner) scanVariable(tok *Token) {
	s.read()
	s.scanIdent(tok)
	tok.Type = variable
}

func (s *Scanner) scanNumber(tok *Token) {
	pos := s.curr
	for isDigit(s.char) {
		s.read()
	}
	if s.char == dot {
		s.read()
		for isDigit(s.char) {
			s.read()
		}
	}
	tok.Literal = string(s.buffer[pos:s.curr])
	tok.Type = number
}

func (s *Scanner) scanText(tok *Token) {
	var (
		buf    []rune
		quote  = s.char
		escape = quote == dquote
	)
	s.read()
	for s.char != quote {
		if escape && s.char == backslash {
			s.read()
			s.char = escapeRune(s.char)
		}
		buf = append(buf, s.char)
		s.read()
	}
	tok.Literal = string(buf)
	tok.Type = text
	s.read()
}

func (s *Scanner) scanPunct(tok *Token) {
	tok.Type = s.char
	if s.char == colon {
		s.read()
		tok.Type = invalid
		if s.char == equal {
			tok.Type = assign
		}
	}
	s.read()
	switch tok.Type {
	case lcurly, lsquare, semicolon, comma:
		s.skipNL()
	default:
	}
}

func (s *Scanner) scanComment(tok *Token) {
	s.read()
	s.read()
	s.skipBlanks()

	pos := s.curr
	for !isNL(s.char) {
		s.read()
	}
	tok.Type = comment
	tok.Literal = string(s.buffer[pos:s.curr])
	s.skipNL()
}

func (s *Scanner) read() {
	if s.curr >= len(s.buffer) {
		s.char = eof
		return
	}
	r, n := utf8.DecodeRune(s.buffer[s.next:])
	s.char, s.curr, s.next = r, s.next, s.next+n
	if s.char == utf8.RuneError {
		s.char = eof
		s.next = len(s.buffer)
	}
}

func (s *Scanner) peek() rune {
	r, _ := utf8.DecodeRune(s.buffer[s.next:])
	return r
}

func (s *Scanner) skipBlanks() {
	for isBlank(s.char) {
		s.read()
	}
}

func (s *Scanner) skipNL() {
	for isNL(s.char) {
		s.read()
	}
}

const (
	eof        = 0
	space      = ' '
	tab        = '\t'
	nl         = '\n'
	cr         = '\r'
	lcurly     = '{'
	rcurly     = '}'
	lsquare    = '['
	rsquare    = ']'
	dollar     = '$'
	equal      = '='
	comma      = ','
	semicolon  = ';'
	dquote     = '"'
	squote     = '\''
	dash       = '-'
	slash      = '/'
	colon      = ':'
	backslash  = '\\'
	dot        = '.'
	underscore = '_'
)

func isNL(r rune) bool {
	return r == nl
}

func isBlank(r rune) bool {
	return r == space || r == tab
}

func isIdent(r rune) bool {
	return isLetter(r) || isDigit(r) || r == underscore || r == dash
}

func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isSpace(r rune) bool {
	return r == space || r == tab
}

func isVariable(r rune) bool {
	return r == dollar
}

func isPunct(r rune) bool {
	return r == lcurly || r == rcurly || r == lsquare || r == rsquare ||
		r == comma || r == semicolon || r == equal || r == colon
}

func isQuote(r rune) bool {
	return r == squote || r == dquote
}

func escapeRune(char rune) rune {
	switch char {
	case 'n':
		char = nl
	case 't':
		char = tab
	case dquote:
	case backslash:
	}
	return char
}

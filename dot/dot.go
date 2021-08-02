package dot

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"unicode/utf8"
)

const (
	kwGraph   = "graph"
	kwInclude = "include"
)

type Graph struct {
	Id    string
	Nodes []*Node
	Attrs map[string][]string
	Edges map[string][]string
}

func createGraph() *Graph {
	return &Graph{
		Attrs: make(map[string][]string),
		Edges: make(map[string][]string),
	}
}

func (g *Graph) insert(n *Node) {

}

type Node struct {
	Id    string
	Attrs map[string][]string
}

func createNode(id string) *Node {
	return &Node{
		Id:    id,
		Attrs: make(map[string][]string),
	}
}

func (n *Node) set(name string, values []string) {
	n.Attrs[name] = values
}

type parser struct {
	frames []*frame
	env    map[string][]string
}

func Parse(r io.Reader) error {
	var p parser
	p.env = make(map[string][]string)
	if err := p.push(r); err != nil {
		return err
	}
	return p.parse()
}

func (p *parser) parse() error {
	var (
		err error
		gph = createGraph()
	)
	for err == nil && !p.done() {
		if p.currIs(comment) {
			p.next()
			continue
		}
		if p.currIs(ident) && p.currIs(lcurly) {
			err = p.unexpectedToken()
			break
		}
		switch {
		case p.peekIs(assign):
			err = p.parseAssignment()
		case p.currLit(kwInclude):
			err = p.parseInclude()
		case p.currLit(kwGraph) || p.currIs(lcurly):
			err = p.parseGraph(gph)
		default:
			err = p.unexpectedToken()
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

	name := p.curr().Literal
	if _, ok := p.env[name]; ok {
		return fmt.Errorf("%s: already defined", name)
	}
	p.next()
	p.next()
	vs, err := p.parseValues()
	if err != nil {
		return err
	}
	p.env[name] = append(p.env[name], vs...)
	fmt.Printf(">> assignment: %s = %s\n", name, p.env[name])
	return p.parseEOL()
}

func (p *parser) parseGraph(g *Graph) error {
	fmt.Println("enter parseGraph")
	defer fmt.Println("leave parseGraph")
	if p.currIs(ident) {
		p.next()
		if p.curr().IsValue() {
			fmt.Printf(">> name: %s\n", p.curr().Literal)
			p.next()
		}
	}
	if !p.currIs(lcurly) {
		return p.unexpectedToken()
	}
	p.next()
	for !p.done() && !p.currIs(rcurly) {
		if p.currIs(comment) {
			p.next()
			continue
		}
		if !p.currIs(ident) {
			return p.unexpectedToken()
		}
		var err error
		if p.peekIs(equal) {
			err = p.parseAttributes(g)
		} else {
			err = p.parseEdge(g)
		}
		if err != nil {
			return err
		}
	}
	if !p.currIs(rcurly) {
		return p.unexpectedToken()
	}
	p.next()
	return p.parseEOL()
}

func (p *parser) parseEdge(g *Graph) error {
	fmt.Println("enter parseEdge")
	defer fmt.Println("leave parseEdge")
	for !p.curr().IsEOL() && !p.currIs(comment) && !p.done() {
		switch p.curr().Type {
		case lcurly:
			p.next()
			for !p.currIs(rcurly) {
				if err := p.parseNode(g); err != nil {
					return err
				}
			}
			if !p.currIs(rcurly) {
				return p.unexpectedToken()
			}
			p.next()
		case ident, text:
			if err := p.parseNode(g); err != nil {
				return err
			}
		default:
			return p.unexpectedToken()
		}
		if p.currIs(edge) {
			if !p.peekIs(ident) && !p.peekIs(lcurly) {
				return p.unexpectedToken()
			}
			p.next()
		}
	}
	return p.parseEOL()
}

func (p *parser) parseNode(g *Graph) error {
	fmt.Println("enter parseNode")
	defer fmt.Println("leave parseNode")

	n := createNode(p.curr().Literal)
	p.next()
	if p.currIs(lsquare) {
		p.next()
		err := p.parseProperties(n)
		if err != nil {
			return err
		}
		if p.currIs(comment) || p.curr().IsEOL() {
			p.next()
			return nil
		}
	}
	fmt.Printf("node: %+v\n", n)
	g.insert(n)
	return nil
}

func (p *parser) parseAttributes(g *Graph) error {
	fmt.Println("enter parseAttributes")
	defer fmt.Println("leave parseAttributes")

	name := p.curr().Literal
	if !p.peekIs(equal) {
		return p.unexpectedToken()
	}
	p.next()
	p.next()
	values, err := p.parseValues()
	if err != nil {
		return err
	}
	if !p.curr().IsEOL() && !p.currIs(comment) {
		return p.unexpectedToken()
	}
	fmt.Printf(">> attributes: %s = %s\n", name, values)
	return p.parseEOL()
}

func (p *parser) parseProperties(n *Node) error {
	fmt.Println("enter parseProperties")
	defer fmt.Println("leave parseProperties")
	for !p.done() && !p.currIs(rsquare) {
		if !p.currIs(ident) && !p.peekIs(equal) {
			return p.unexpectedToken()
		}
		name := p.curr().Literal
		p.next()
		p.next()
		values, err := p.parseValues()
		if err != nil {
			return err
		}
		if !p.currIs(comma) && !p.currIs(rsquare) {
			return p.unexpectedToken()
		}
		if p.currIs(comma) {
			p.next()
		}
		if p.currIs(comment) {
			p.next()
		}
		n.set(name, values)
		fmt.Printf(">> properties: %s = %s\n", name, values)
	}
	if !p.currIs(rsquare) {
		return p.unexpectedToken()
	}
	p.next()
	return nil
}

func (p *parser) parseValues() ([]string, error) {
	var vs []string
	for p.curr().IsValue() {
		if p.currIs(variable) {
			xs, ok := p.env[p.curr().Literal]
			if !ok {
				return nil, fmt.Errorf("%s: undefined variable", p.curr().Literal)
			}
			vs = append(vs, xs...)
		} else {
			vs = append(vs, p.curr().Literal)
		}
		p.next()
	}
	return vs, nil
}

func (p *parser) parseEOL() error {
	if p.currIs(comment) {
		p.next()
		return nil
	}
	if !p.curr().IsEOL() {
		return p.unexpectedToken()
	}
	p.next()
	return nil
}

func (p *parser) parseInclude() error {
	fmt.Println("enter parseInclude")
	defer fmt.Println("leave parseInclude")
	p.next()
	if !p.currIs(text) && !p.currIs(ident) {
		return p.unexpectedToken()
	}
	file := p.curr().Literal
	fmt.Printf(">> include: %s\n", p.curr().Literal)
	p.next()
	if err := p.parseEOL(); err != nil {
		return err
	}
	r, err := os.Open(file)
	if err != nil {
		return err
	}
	defer r.Close()
	return p.push(r)
}

func (p *parser) unexpectedToken() error {
	return unexpectedToken(p.curr())
}

func (p *parser) next() {
	z := len(p.frames)
	if z == 0 {
		return
	}
	z--
	p.frames[z].next()
	if p.frames[z].done() {
		p.pop()
	}
}

func (p *parser) done() bool {
	return len(p.frames) == 0
}

func (p *parser) curr() Token {
	z := len(p.frames)
	if z == 0 {
		return Token{}
	}
	return p.frames[z-1].curr
}

func (p *parser) currIs(k rune) bool {
	return p.curr().Type == k
}

func (p *parser) currLit(str string) bool {
	return p.curr().Literal == str
}

func (p *parser) peek() Token {
	z := len(p.frames)
	if z == 0 {
		return Token{}
	}
	return p.frames[z-1].peek
}

func (p *parser) peekIs(k rune) bool {
	return p.peek().Type == k
}

func (p *parser) peekLit(str string) bool {
	return p.peek().Literal == str
}

func (p *parser) push(r io.Reader) error {
	f, err := createFrame(r)
	if err == nil {
		p.frames = append(p.frames, f)
	}
	return err
}

func (p *parser) pop() {
	z := len(p.frames)
	if z == 0 {
		return
	}
	p.frames = p.frames[:z-1]
}

type frame struct {
	scan *Scanner
	curr Token
	peek Token
}

func createFrame(r io.Reader) (*frame, error) {
	s, err := Scan(r)
	if err != nil {
		return nil, err
	}
	f := frame{
		scan: s,
	}
	f.next()
	f.next()
	return &f, nil
}

func (f *frame) next() {
	f.curr = f.peek
	f.peek = f.scan.Scan()
}

func (f *frame) done() bool {
	return f.curr.Type == eof
}

const (
	ident rune = -(iota + 1)
	text
	number
	boolean
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
	switch t.Type {
	case number, text, ident, variable, boolean:
		return true
	default:
		return false
	}
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
	case boolean:
		prefix = "boolean"
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

const (
	kwTrue  = "true"
	kwFalse = "false"
)

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
	if tok.Literal == kwTrue || tok.Literal == kwFalse {
		tok.Type = boolean
	}
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

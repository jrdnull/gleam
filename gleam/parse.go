package gleam

import (
	"errors"
	"fmt"
)

type nodeType int

const (
	nodeSymbol nodeType = iota
	nodeNumber
	nodeList
)

type node interface {
	String() string
}

type listNode struct {
	nodes []node
}

func (n *listNode) String() string {
	s := ""
	for i, c := range n.nodes {
		if i > 0 {
			s += " "
		}
		s += c.String()
	}
	return "(" + s + ")"
}

type symbolNode struct {
	val string
}

func (n *symbolNode) String() string {
	return n.val
}

type numberNode struct {
	val string
}

func (n *numberNode) String() string {
	return n.val
}

type parser struct {
	tokens []token
	pos    int
	err    error
	node   node
}

func parse(tokens []token) (node, error) {
	parser := parser{
		tokens: tokens,
		pos:    -1,
	}
	fmt.Printf("%#v\n", tokens) // output for debug

	node, err := parser.parse()
	if err != nil {
		return nil, err
	}

	return node, nil
}

func (p *parser) parse() (node, error) {
	t := p.next()
	switch t.typ {
	case tokenLeftParen:
		l := listNode{}
		for {
			next := p.peek().typ
			if next == tokenEOF {
				return nil, errors.New("unclosed left paren")
			} else if next == tokenRightParen {
				break
			}

			n, err := p.parse()
			if err != nil {
				return nil, err
			}

			l.nodes = append(l.nodes, n)
		}
		return &l, nil
	case tokenRightParen:
		return nil, errors.New("unexpected right paren")
	case tokenNumber:
		return &numberNode{val: t.val}, nil
	case tokenSymbol:
		return &symbolNode{val: t.val}, nil
	default:
		return nil, fmt.Errorf("unexpected token: %#v", t)
	}
}

func (p *parser) next() token {
	p.pos++
	return p.tokens[p.pos]
}

func (p *parser) peek() token {
	return p.tokens[p.pos+1]
}

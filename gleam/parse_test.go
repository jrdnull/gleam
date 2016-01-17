package gleam

import (
	"reflect"
	"testing"

	"github.com/kr/pretty"
)

var parseTests = []struct {
	test string
	in   []token
	out  node
}{
	{"(+ (+ 1 2) (+ 3 3))", []token{
		token{tokenLeftParen, "("},
		token{tokenSymbol, "+"},
		token{tokenLeftParen, "("},
		token{tokenSymbol, "+"},
		token{tokenNumber, "1"},
		token{tokenNumber, "2"},
		token{tokenRightParen, ")"},
		token{tokenLeftParen, "("},
		token{tokenSymbol, "+"},
		token{tokenNumber, "3"},
		token{tokenNumber, "3"},
		token{tokenRightParen, ")"},
		token{tokenRightParen, ")"},
		token{tokenEOF, ""},
	}, &listNode{[]node{
		&symbolNode{"+"},
		&listNode{[]node{
			&symbolNode{"+"},
			&numberNode{1},
			&numberNode{2},
		}},
		&listNode{[]node{
			&symbolNode{"+"},
			&numberNode{3},
			&numberNode{3},
		}},
	}}},
}

func TestParse(t *testing.T) {
	for _, pt := range parseTests {
		n, err := parse(pt.in)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(n, pt.out) {
			t.Errorf("case: %s", pt.test)
			t.Errorf(pretty.Sprintf("got: %# v\n", n))
			t.Errorf(pretty.Sprintf("want: % #v\n", pt.out))
		}
	}
}

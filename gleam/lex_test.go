package gleam

import "testing"

var lexTests = []struct {
	in  string
	out []item
}{
	{"()", []item{
		item{itemLeftParen, "("},
		item{itemRightParen, ")"},
		item{itemEOF, ""},
	}},
	{"(", []item{
		item{itemLeftParen, "("},
		item{itemError, "unclosed left paren"},
	}},
	{")", []item{
		item{itemError, "unexpected right paren"},
	}},
	{"123", []item{
		item{itemNumber, "123"},
		item{itemEOF, ""},
	}},
	{"123", []item{
		item{itemNumber, "123"},
		item{itemEOF, ""},
	}},
	{"+123", []item{
		item{itemNumber, "+123"},
		item{itemEOF, ""},
	}},
	{"-123", []item{
		item{itemNumber, "-123"},
		item{itemEOF, ""},
	}},
	{"123.456", []item{
		item{itemNumber, "123.456"},
		item{itemEOF, ""},
	}},
	{"foo", []item{
		item{itemSymbol, "foo"},
		item{itemEOF, ""},
	}},
	{"+ 1234", []item{
		item{itemSymbol, "+"},
		item{itemNumber, "1234"},
		item{itemEOF, ""},
	}},
	{"(* 2 (+ 5 7))", []item{
		item{itemLeftParen, "("},
		item{itemSymbol, "*"},
		item{itemNumber, "2"},
		item{itemLeftParen, "("},
		item{itemSymbol, "+"},
		item{itemNumber, "5"},
		item{itemNumber, "7"},
		item{itemRightParen, ")"},
		item{itemRightParen, ")"},
		item{itemEOF, ""},
	}},
}

func TestLex(t *testing.T) {
	for _, lt := range lexTests {
		_, items := lex(lt.in)

		var tokens []item
		for i := range items {
			tokens = append(tokens, i)
		}

		if len(tokens) != len(lt.out) {
			t.Errorf("lex(%q) = %v, want %v", lt.in, tokens, lt.out)
			continue
		}

		for i := range tokens {
			if tokens[i] != lt.out[i] {
				t.Errorf("lex(%q) = %#v, want %#v", lt.in, tokens, lt.out)
				continue
			}
		}
	}
}

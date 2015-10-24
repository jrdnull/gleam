package gleam

import "testing"

var lexTests = []struct {
	in  string
	out []token
}{
	{"()", []token{
		token{tokenLeftParen, "("},
		token{tokenRightParen, ")"},
		token{tokenEOF, ""},
	}},
	{"(", []token{
		token{tokenLeftParen, "("},
		token{tokenEOF, ""},
	}},
	{")", []token{
		token{tokenRightParen, ")"},
		token{tokenEOF, ""},
	}},
	{"123", []token{
		token{tokenNumber, "123"},
		token{tokenEOF, ""},
	}},
	{"123", []token{
		token{tokenNumber, "123"},
		token{tokenEOF, ""},
	}},
	{"+123", []token{
		token{tokenNumber, "+123"},
		token{tokenEOF, ""},
	}},
	{"-123", []token{
		token{tokenNumber, "-123"},
		token{tokenEOF, ""},
	}},
	{"123.456", []token{
		token{tokenNumber, "123.456"},
		token{tokenEOF, ""},
	}},
	{"foo", []token{
		token{tokenSymbol, "foo"},
		token{tokenEOF, ""},
	}},
	{"+ 1234", []token{
		token{tokenSymbol, "+"},
		token{tokenNumber, "1234"},
		token{tokenEOF, ""},
	}},
	{"(* 2 (+ 5 7))", []token{
		token{tokenLeftParen, "("},
		token{tokenSymbol, "*"},
		token{tokenNumber, "2"},
		token{tokenLeftParen, "("},
		token{tokenSymbol, "+"},
		token{tokenNumber, "5"},
		token{tokenNumber, "7"},
		token{tokenRightParen, ")"},
		token{tokenRightParen, ")"},
		token{tokenEOF, ""},
	}},
	{"(foo bar)", []token{
		token{tokenLeftParen, "("},
		token{tokenSymbol, "foo"},
		token{tokenSymbol, "bar"},
		token{tokenRightParen, ")"},
		token{tokenEOF, ""},
	}},
	{"(foo bar baz)", []token{
		token{tokenLeftParen, "("},
		token{tokenSymbol, "foo"},
		token{tokenSymbol, "bar"},
		token{tokenSymbol, "baz"},
		token{tokenRightParen, ")"},
		token{tokenEOF, ""},
	}},
	{"(foo bar (baz qux 7))", []token{
		token{tokenLeftParen, "("},
		token{tokenSymbol, "foo"},
		token{tokenSymbol, "bar"},
		token{tokenLeftParen, "("},
		token{tokenSymbol, "baz"},
		token{tokenSymbol, "qux"},
		token{tokenNumber, "7"},
		token{tokenRightParen, ")"},
		token{tokenRightParen, ")"},
		token{tokenEOF, ""},
	}},
}

func TestLex(t *testing.T) {
	for _, lt := range lexTests {
		tokens := lex(lt.in)
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

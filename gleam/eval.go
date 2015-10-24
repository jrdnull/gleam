package gleam

import "fmt"

func Eval(input string) string {
	tokens, err := parse(lex(input))
	fmt.Printf("%#v\n", tokens) // output for debug
	if err != nil {
		return err.Error()
	}

	return tokens.String()
}

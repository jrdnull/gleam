package gleam

import (
	"errors"
	"fmt"

	"github.com/kr/pretty"
)

type fun func(args ...node) (node, error)

var env = map[string]fun{
	"+": func(args ...node) (node, error) {
		return reduceNumbers(args, func(acc, n float64) float64 {
			return acc + n
		}, 0)
	},
	"-": func(args ...node) (node, error) {
		if len(args) == 0 {
			return nil, errors.New("requires at least 1 argument")
		} else if len(args) == 1 {
			return &numberNode{-args[0].(*numberNode).val}, nil
		}

		return reduceNumbers(args, func(acc, n float64) float64 {
			return acc - n
		}, args[0].(*numberNode).val)
	},
	"*": func(args ...node) (node, error) {
		return reduceNumbers(args, func(acc, n float64) float64 {
			return acc * n
		}, 1)
	},
	"/": func(args ...node) (node, error) {
		if len(args) == 0 {
			return nil, errors.New("requires at least 1 argument")
		}
		return reduceNumbers(args[1:], func(acc, n float64) float64 {
			return acc / n
		}, args[0].(*numberNode).val)
	},
}

var userEnv = make(map[string]node)

func reduceNumbers(ns []node, f func(acc, n float64) float64, initial float64) (node, error) {
	acc := initial
	for _, n := range ns {
		n, ok := n.(*numberNode)
		if !ok {
			return nil, errors.New("expected number")
		}
		acc = f(acc, n.val)
	}
	return &numberNode{acc}, nil
}

func Eval(input string) string {
	root, err := parse(lex(input))
	if err != nil {
		return err.Error()
	}

	fmt.Println("ast:")
	pretty.Println(root)

	n, err := eval(root)
	if err != nil {
		return "error: " + err.Error()
	}
	return "result:\n" + n.String()
}

func eval(n node) (node, error) {
	switch n := n.(type) {
	case *listNode:
		if len(n.nodes) == 0 {
			return n, nil
		}

		if s, ok := n.nodes[0].(*symbolNode); ok {
			if s.val == "define" {
				if len(n.nodes) != 3 {
					return nil, errors.New("expected identifier and expression")
				}
				userEnv[n.nodes[1].String()] = n.nodes[2]
				return n, nil
			}

			f, ok := env[s.val]
			if !ok {
				return nil, fmt.Errorf("%s undefined", s.val)
			}
			var args []node
			for _, childNode := range n.nodes[1:] {
				evald, err := eval(childNode)
				if err != nil {
					return nil, err
				}
				args = append(args, evald)
			}
			return f(args...)
		}
		return n, nil
	case *symbolNode:
		if n, ok := userEnv[n.val]; ok {
			return n, nil
		}
		return nil, fmt.Errorf("%s undefined", n.val)
	default:
		return n, nil
	}
}

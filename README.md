# gleam
lisp interpreter written in Go, lexer based on the one presented by Rob Pike in his talk [Lexical Scanning in Go](https://www.youtube.com/watch?v=HxaD_trXwRE).

## Issues
### Lexer
- Lex number without a space or ) before symbol should be invalid

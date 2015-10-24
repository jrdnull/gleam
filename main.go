package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jrdnull/gleam/gleam"
)

func main() {
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("gleam> ")
		line, _ := r.ReadString('\n')
		fmt.Println(gleam.Eval(line))
	}
}

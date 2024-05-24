package main

import (
	"fmt"
	"monkeylang/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()

	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Monkey Programming Language!\n", user.Username)
	fmt.Print("Reading from scratch/input.mkl ...\n\n")

	//fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}

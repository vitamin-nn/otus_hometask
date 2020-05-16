package main

import (
	"errors"
	"fmt"
	"os"
)

const (
	minArgs            = 3
	defaultExitErrCode = 111
)

var (
	ErrInvalidArgumentCount = errors.New("invalid argument count")
)

func main() {
	args := os.Args
	if len(args) < minArgs {
		fmt.Println(ErrInvalidArgumentCount)
		os.Exit(defaultExitErrCode)
	}
	d := args[1]
	child := args[2:]

	env, err := ReadDir(d)
	if err != nil {
		fmt.Println(err)
		os.Exit(defaultExitErrCode)
	}
	code := RunCmd(child, env)
	os.Exit(code)
}

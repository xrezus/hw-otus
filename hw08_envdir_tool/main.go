package main

import (
	"errors"
	"log"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		log.Fatal("Not enough arguments. Usage: go-envdir /path/to/evndir command arg1 arg2...")
	}

	env, err := ReadDir(args[1])
	if err != nil {
		if !errors.Is(err, ErrWrongFileName) {
			log.Fatal(err)
		}
		log.Fatal(err)
	}

	exitCode := RunCmd(args[2:], env)
	os.Exit(exitCode)
}

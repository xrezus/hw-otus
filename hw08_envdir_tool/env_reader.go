package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var ErrWrongFileName = fmt.Errorf("some env files have wrong names: include '=', skipped")

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	envList := make(Environment)
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var fN string
	var enV EnvValue

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fN = file.Name()

		if strings.Contains(fN, "=") {
			err = ErrWrongFileName
			continue
		}

		readFile, err := os.Open(filepath.Join(dir, fN))
		defer func() { readFile.Close() }()
		if err != nil {
			return nil, err
		}

		rd := bufio.NewReader(readFile)
		s, err := rd.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		}

		if len(s) == 0 {
			enV.NeedRemove = true
		}

		s = strings.TrimRight(s, " \t\n")
		enV.Value = string(bytes.ReplaceAll([]byte(s), []byte{0}, []byte{10}))
		envList[fN] = enV
	}

	return envList, err
}

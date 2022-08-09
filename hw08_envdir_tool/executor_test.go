package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	env := Environment{"BAR": {"bar", false}, "UNSET": {"", true}}
	cmdOk := []string{"/bin/bash", "./testdata/echo.sh", "arg1=1", "arg2=2"}
	cmdFalse := []string{"/bin/bas", "./testdata/echo.sh", "arg1=1", "arg2=2"}
	t.Run("os variables updated correctly", func(t *testing.T) {
		osExit := RunCmd(cmdOk, env)
		v := os.Getenv("BAR")

		require.Equal(t, 0, osExit)
		require.Equal(t, "bar", v)
	})
	t.Run("os variables updated correctly and unset", func(t *testing.T) {
		osExit := RunCmd(cmdOk, env)
		v := os.Getenv("UNSET")

		require.Equal(t, 0, osExit)
		require.Equal(t, "", v)
	})
	t.Run("no such file or directory returns 1", func(t *testing.T) {
		osExit := RunCmd(cmdFalse, env)
		require.Equal(t, 1, osExit)
	})
	t.Run("os local variables updated correctly", func(t *testing.T) {
		f1, err := os.CreateTemp("./testdata/env", "example*.sh")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(f1.Name())
		f2, err := os.CreateTemp("./testdata/env", "result")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(f2.Name())
		if _, err := f1.Write([]byte("#!/usr/bin/env bash\necho $*>" + f2.Name())); err != nil {
			log.Fatal(err)
		}
		cmdOk = []string{"/bin/bash", f1.Name(), "arg1=1"}
		osExit := RunCmd(cmdOk, env)
		rd := bufio.NewReader(f2)
		v, _ := rd.ReadString('\n')
		v = strings.TrimRight(v, " \t\n")

		require.Equal(t, 0, osExit)
		require.Equal(t, "arg1=1", v)
	})
}

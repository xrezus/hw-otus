package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for enVar, value := range env {
		if value.NeedRemove {
			os.Unsetenv(enVar)
			continue
		}

		if _, ok := os.LookupEnv(enVar); ok {
			os.Unsetenv(enVar)
		}

		os.Setenv(enVar, value.Value)
	}

	name, args := cmd[0], cmd[1:]

	proc := exec.Command(name, args...)

	proc.Env = append(os.Environ(), args...)

	proc.Stdout = os.Stdout
	proc.Stderr = os.Stderr
	proc.Stdin = os.Stdin

	if err := proc.Run(); err != nil {
		return 1
	}

	return
}

package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("ErrWrongFileName filename include '='.", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("./", "example=.")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(tmpFile.Name())

		_, err = ReadDir("./")

		require.Truef(t, errors.Is(err, ErrWrongFileName), "actual err - %v", err)
	})
	t.Run("environment variables received correctly", func(t *testing.T) {
		envName := "BAR"
		dir := "./testdata/env"
		v2 := "bar"

		env, err := ReadDir(dir)

		v1, ok1 := env[envName]

		require.Truef(t, errors.Is(err, nil), "actual err - %v", err)
		require.Equal(t, true, ok1)
		require.Equal(t, v2, v1.Value)
		require.Equal(t, false, v1.NeedRemove)
	})
	t.Run("empty file content - need to remove variable", func(t *testing.T) {
		envName := "UNSET"
		dir := "./testdata/env"

		env, err := ReadDir(dir)

		v1, ok1 := env[envName]
		fmt.Println(env)

		require.Truef(t, errors.Is(err, nil), "actual err - %v", err)
		require.Equal(t, true, ok1)
		require.Equal(t, true, v1.NeedRemove)
	})
}

package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		// uncomment if task with asterisk completed
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},

		{input: `da😀3qe`, expected: `da😀😀😀qe`},
		{input: `ПриветМир3!!!`, expected: `ПриветМиррр!!!`},
		{input: `При0ф`, expected: `Прф`},
		{input: `spa 3ce`, expected: `spa   ce`},
		{input: `Two	2tab`, expected: `Two		tab`},
		{input: `a0b0c0d0e0f0`, expected: ``},
		{input: `\1`, expected: `1`},
		{input: `\\0`, expected: ``},
		{input: `\0\0\0`, expected: `000`},
		{input: `a3+1`, expected: `aaa+`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{
		"3abc",
		"45",
		"aaa10b",
		"3\a",
		"3+1",
	}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}

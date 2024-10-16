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
		err      error
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

		// Строка только с экранированными символами
		{input: `\a`, expected: "", err: ErrInvalidString},
		{input: `\\3a`, expected: `\\\a`, err: nil},
		{input: `\\\\`, expected: `\\`, err: nil},

		// Строка с нулевыми повторами
		{input: "a0b0c0", expected: "", err: nil},

		// Несколько экранированных символов
		{input: `qwe\3\4`, expected: "qwe34", err: nil},

		// Экранирование цифр и слешей
		{input: `\1a\2b\3`, expected: "1a2b3", err: nil},

		// Неправильное экранирование
		{input: `a\`, expected: "", err: ErrInvalidString},
		{input: `a\g`, expected: "", err: ErrInvalidString},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			if tc.err != nil {
				require.Error(t, err)
				require.Truef(t, errors.Is(err, tc.err), "expected error %q, got %q", tc.err, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}

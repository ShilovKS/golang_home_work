package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const output = "./testdata/output.txt"

func cleanup() {
	_ = os.Remove(output)
}

func TestCopy0_0(t *testing.T) {
	defer cleanup()
	err := Copy("testdata/input.txt", output, 0, 0)
	require.NoError(t, err)
	fact, err := os.ReadFile(output)
	require.NoError(t, err)
	expected, err := os.ReadFile("testdata/out_offset0_limit0.txt")
	require.NoError(t, err)
	require.Equal(t, expected, fact)
}

func TestCopy0_10(t *testing.T) {
	defer cleanup()
	err := Copy("testdata/input.txt", output, 0, 10)
	require.NoError(t, err)
	fact, err := os.ReadFile(output)
	require.NoError(t, err)
	expected, err := os.ReadFile("testdata/out_offset0_limit10.txt")
	require.NoError(t, err)
	require.Equal(t, expected, fact)
}

func TestCopy0_1000(t *testing.T) {
	defer cleanup()
	err := Copy("testdata/input.txt", output, 0, 1000)
	require.NoError(t, err)
	fact, err := os.ReadFile(output)
	require.NoError(t, err)
	expected, err := os.ReadFile("testdata/out_offset0_limit1000.txt")
	require.NoError(t, err)
	require.Equal(t, expected, fact)
}

func TestCopy0_10000(t *testing.T) {
	defer cleanup()
	err := Copy("testdata/input.txt", output, 0, 10000)
	require.NoError(t, err)
	fact, err := os.ReadFile(output)
	require.NoError(t, err)
	expected, err := os.ReadFile("testdata/out_offset0_limit10000.txt")
	require.NoError(t, err)
	require.Equal(t, expected, fact)
}

func TestCopy100_1000(t *testing.T) {
	defer cleanup()
	err := Copy("testdata/input.txt", output, 100, 1000)
	require.NoError(t, err)
	fact, err := os.ReadFile(output)
	require.NoError(t, err)
	expected, err := os.ReadFile("testdata/out_offset100_limit1000.txt")
	require.NoError(t, err)
	require.Equal(t, expected, fact)
}

func TestCopy6000_1000(t *testing.T) {
	defer cleanup()
	err := Copy("testdata/input.txt", output, 6000, 1000)
	require.NoError(t, err)
	expected, err := os.ReadFile("testdata/out_offset6000_limit1000.txt")
	require.NoError(t, err)
	fact, err := os.ReadFile(output)
	require.NoError(t, err)
	require.Equal(t, expected, fact)
}

func TestCopyOverlap(t *testing.T) {
	defer cleanup()
	err := Copy("testdata/input.txt", "testdata/input.txt", 0, 0)
	require.ErrorIs(t, err, ErrFileOverlap)
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(Environment)

	dirs, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	for _, currentDir := range dirs {
		if currentDir.IsDir() || currentDir.Name()[0] == '.' {
			continue
		}
		filePath := filepath.Join(dir, currentDir.Name())
		file, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		var line string
		if scanner.Scan() {
			line = scanner.Text()
			line = strings.TrimRight(line, " \t\n\r")
			line = strings.ReplaceAll(line, "\x00", "\n")
		}

		if line == "" {
			env[currentDir.Name()] = EnvValue{NeedRemove: true}
		} else {
			env[currentDir.Name()] = EnvValue{Value: line, NeedRemove: false}
		}

		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("error reading file %s: %w", filePath, err)
		}
	}
	return env, nil
}

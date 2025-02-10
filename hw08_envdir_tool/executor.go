package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for key, val := range env {
		if val.NeedRemove {
			err := os.Unsetenv(key)
			if err != nil {
				return 1
			}
			continue
		}

		err := os.Setenv(key, val.Value)
		if err != nil {
			return 1
		}
	}

	if len(cmd) == 0 {
		return 1
	}
	executable, err := exec.LookPath(cmd[0])
	if err != nil {
		return 1
	}

	command := exec.Command(executable, cmd[1:]...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err = command.Run()
	if err != nil {
		return command.ProcessState.ExitCode()
	}
	return 0
}

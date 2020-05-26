package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	var args []string
	if len(cmd) > 1 {
		args = cmd[1:]
	}
	command := cmd[0]
	execCmd := exec.Command(command, args...)
	execCmd.Env = prepareEnv(env)
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	err := execCmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode()
		}
		return defaultExitErrCode
	}
	return 0
}

func prepareEnv(env Environment) []string {
	envList := os.Environ()
	for k, v := range env {
		envList = append(envList, fmt.Sprintf("%s=%s", k, v))
	}
	return envList
}

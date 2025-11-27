package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...)
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout
	command.Env = os.Environ()
	for k, v := range env {
		if !v.NeedRemove {
			command.Env = append(command.Env, fmt.Sprintf("%s=%s", k, v.Value))
		} else {
			for i := range command.Env {
				if strings.HasPrefix(command.Env[i], fmt.Sprintf("%s=", k)) {
					command.Env = append(command.Env[:i], command.Env[i+1:]...)
					break
				}
			}
		}
	}
	err := command.Run()
	if err != nil {
		return -1
	}
	return 0
}

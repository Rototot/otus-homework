package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		log.Println("empty command")
		return 1
	}

	var cmdEnvs = mergeWithGlobal(env)

	// init command
	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	command.Env = cmdEnvs
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err := command.Run()
	if err != nil {
		log.Fatal(err)
	}

	return 0
}

func mergeWithGlobal(envs Environment) []string {
	var sourceEnvs = envs
	var resultEnvs = make([]string, 0, len(envs)+len(os.Environ()))

	for _, value := range os.Environ() {
		envParts := strings.Split(value, "=")

		// if exists in envs, which grabbed from envdir
		sourceValue, ok := sourceEnvs[envParts[0]]
		if ok {
			resultEnvs = append(resultEnvs, envParts[0]+"="+sourceValue)
			delete(sourceEnvs, envParts[0])
		}
	}

	// append another envs
	for k, v := range sourceEnvs {
		resultEnvs = append(resultEnvs, k+"="+v)
	}

	return resultEnvs
}

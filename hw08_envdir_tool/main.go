package main

import (
	"log"
	"os"
)

func main() {
	var args = os.Args
	if len(args) <= 2 {
		log.Fatalln("incorrect qty arguments. Please set by template: <environment dir> <command> [arguments] ")
	}

	var envDir = args[1]
	dirEnvs, err := ReadDir(envDir)
	if err != nil {
		log.Fatalln(err)
	}

	code := RunCmd(args[2:], dirEnvs)

	os.Exit(code)
}

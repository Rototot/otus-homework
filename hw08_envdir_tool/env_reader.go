package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
)

type Environment map[string]string

type rawEnvValue struct {
	key   string
	value string
	err   error
}
type filter func(v string) string

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	rawEnvs, err := grabRawEnvs(dir)
	if err != nil {
		return nil, err
	}

	var envs = make(Environment, cap(rawEnvs))
	for rawEnv := range rawEnvs {
		key, err := extractEnvKey(rawEnv)
		if err != nil {
			return nil, err
		}

		value, err := extractEnvValue(rawEnv)
		if err != nil {
			return nil, err
		}

		if key != "" {
			envs[key] = value
		}
	}

	return envs, err
}

func extractEnvKey(raw rawEnvValue) (string, error) {
	if raw.err != nil {
		return "", raw.err
	}

	keyFilters := []filter{
		strings.TrimSpace,
		// check forbidden symbols
		func(v string) string {
			if strings.Contains(v, "=") {
				return ""
			}

			return v
		},
	}
	key := runFilters(raw.key, keyFilters)

	return key, nil
}

func extractEnvValue(raw rawEnvValue) (string, error) {
	if raw.err != nil {
		return "", raw.err
	}

	valueFilters := []filter{
		func(v string) string { return strings.TrimRight(v, "\n\t") },
		func(v string) string { return strings.ReplaceAll(v, string(0x00), "\n") },
	}
	value := runFilters(raw.value, valueFilters)

	return value, nil
}

func runFilters(value string, filters []filter) string {
	for _, filter := range filters {
		value = filter(value)
	}

	return value
}

func grabRawEnvs(dir string) (<-chan rawEnvValue, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var foundedEnvs = make(chan rawEnvValue, len(files))
	go func() {
		defer close(foundedEnvs)

		var wg sync.WaitGroup
		for _, file := range files {
			if !file.IsDir() {
				wg.Add(1)
				go func(fileInfo os.FileInfo) {
					var result rawEnvValue

					defer wg.Done()
					defer func() { foundedEnvs <- result }()

					file, err := os.Open(path.Join(dir, fileInfo.Name()))
					if err != nil {
						result.err = err
						return
					}
					defer file.Close()

					buf := bufio.NewReader(file)
					line, _, err := buf.ReadLine()
					if err != nil && err != io.EOF {
						result.err = err
						return
					}

					result.key = fileInfo.Name()
					result.value = string(line)
				}(file)
			}
		}
		wg.Wait()
	}()

	return foundedEnvs, nil
}

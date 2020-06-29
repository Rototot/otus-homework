package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestReadDir(t *testing.T) {
	pwd, err := os.Getwd()
	assert.NoError(t, err)

	t.Run("simple", func(t *testing.T) {
		var envDir = path.Join(pwd, "testdata/env")
		//var expectedErr error = nil
		var expectedEnvs Environment = map[string]string{
			"FOO":   "   foo\nwith new line",
			"BAR":   "bar",
			"HELLO": "\"hello\"",
			"UNSET": "",
		}
		envs, err := ReadDir(envDir)

		assert.NoError(t, err)
		assert.EqualValues(t, expectedEnvs, envs)
	})

	t.Run("when directory does not exists", func(t *testing.T) {
		var envDir = path.Join(pwd, "testdata/not_exists")
		envs, err := ReadDir(envDir)

		assert.EqualError(t, err, fmt.Sprintf("open %s: no such file or directory", envDir))
		assert.Nil(t, envs)
	})

	t.Run("when invalid file name", func(t *testing.T) {
		var envDir = path.Join(pwd, "testdata/env_invalid")
		//var expectedErr error = nil
		var expectedEnvs Environment = map[string]string{
			"FOO": "   foo\nwith new line",
		}
		envs, err := ReadDir(envDir)

		assert.NoError(t, err)
		assert.EqualValues(t, expectedEnvs, envs)
	})
}

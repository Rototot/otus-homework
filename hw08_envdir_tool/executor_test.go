package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRunCmd(t *testing.T) {
	type args struct {
		cmd []string
		env Environment
	}

	t.Run("when is ok", func(t *testing.T) {
		var expectedCode int = 0
		var args = args{
			cmd: []string{
				"ls",
				"-la",
			},
			env: map[string]string{
				"FOO": "BAR",
			},
		}

		result := RunCmd(args.cmd, args.env)
		assert.Equal(t, expectedCode, result)
	})

	t.Run("when empty command", func(t *testing.T) {
		var expectedCode int = 1
		var args = args{
			env: map[string]string{
				"FOO": "BAR",
			},
		}

		result := RunCmd(args.cmd, args.env)
		assert.Equal(t, expectedCode, result)
	})

	t.Run("when empty env", func(t *testing.T) {
		var expectedCode int = 0
		var args = args{
			cmd: []string{
				"ls",
				"-la",
			},
		}

		result := RunCmd(args.cmd, args.env)
		assert.Equal(t, expectedCode, result)
	})
}

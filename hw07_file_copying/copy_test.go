package main

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCopy(t *testing.T) {
	type copyCase struct {
		name       string
		args       copyArgs
		wantErr    error
		expectFile string
	}

	copyCases := []copyCase{
		{
			name: "simple",
			args: copyArgs{
				fromPath: "testdata/input.txt",
				toPath:   "out.txt",
				offset:   0,
				limit:    0,
			},
			wantErr:    nil,
			expectFile: "testdata/input.txt",
		},
		{
			name: "when set limit",
			args: copyArgs{
				fromPath: "testdata/input.txt",
				toPath:   "out_limit_10.txt",
				offset:   0,
				limit:    10,
			},
			wantErr:    nil,
			expectFile: "testdata/out_offset0_limit10.txt",
		},
		{
			name: "when set limit and offset 1",
			args: copyArgs{
				fromPath: "testdata/input.txt",
				toPath:   "out_limit_10.txt",
				offset:   100,
				limit:    1000,
			},
			wantErr:    nil,
			expectFile: "testdata/out_offset100_limit1000.txt",
		},
		{
			name: "when set limit 100 and offset 1000",
			args: copyArgs{
				fromPath: "testdata/input.txt",
				toPath:   "out_limit_10.txt",
				offset:   100,
				limit:    1000,
			},
			wantErr:    nil,
			expectFile: "testdata/out_offset100_limit1000.txt",
		},
		{
			name: "when set limit 1000 and offset 6000",
			args: copyArgs{
				fromPath: "testdata/input.txt",
				toPath:   "out_limit_10.txt",
				offset:   6000,
				limit:    1000,
			},
			wantErr:    nil,
			expectFile: "testdata/out_offset6000_limit1000.txt",
		},
	}

	// cases
	for _, tt := range copyCases {
		t.Run("case: " + tt.name, func(t *testing.T) {
			expectFile, err := os.Open(tt.expectFile)
			assert.NoError(t, err)
			defer expectFile.Close()

			expectedFileInfo, err := expectFile.Stat()
			assert.NoError(t, err)

			result := Copy(tt.args.fromPath, tt.args.toPath, tt.args.offset, tt.args.limit)
			assert.NoError(t, result)
			assert.FileExists(t, tt.args.toPath)

			actualFile, err := os.Open(tt.args.toPath)
			assert.NoError(t, err)

			actualFileInfo, err := actualFile.Stat()
			assert.NoError(t, err)

			assert.Equal(t, expectedFileInfo.Size(), actualFileInfo.Size())

			if err := os.Remove(tt.args.toPath); err != nil {
				t.Fatal(err)
			}
		})
	}

	// validate
	validateCases := []copyCase{
		{
			name: "limit < 0",
			args: copyArgs{
				fromPath: "testdata/input.txt",
				toPath:   "out.txt",
				offset:   -2,
				limit:    0,
			},
			wantErr: ErrInvalidArgument,
		},
		{
			name: "offset < 0",
			args: copyArgs{
				fromPath: "testdata/input.txt",
				toPath:   "out.txt",
				offset:   0,
				limit:    -1,
			},
			wantErr: ErrInvalidArgument,
		},
		{
			name: "empty source path",
			args: copyArgs{
				fromPath: "",
				toPath:   "out.txt",
				offset:   0,
				limit:    0,
			},
			wantErr: ErrEmptyPath,
		},
		{
			name: "empty dst path",
			args: copyArgs{
				fromPath: "testdata/input.txt",
				toPath:   "",
				offset:   0,
				limit:    0,
			},
			wantErr: ErrEmptyPath,
		},
		{
			name: "offset more than source size",
			args: copyArgs{
				fromPath: "testdata/input.txt",
				toPath:   "out.txt",
				offset:   1024 * 1024,
				limit:    0,
			},
			wantErr: ErrOffsetExceedsFileSize,
		},
		{
			name: "source path is directory",
			args: copyArgs{
				fromPath: "testdata",
				toPath:   "out.txt",
				offset:   0,
				limit:    0,
			},
			wantErr: errors.New("testdata is not a file"),
		},
		{
			name: "file with unknown length",
			args: copyArgs{
				fromPath: "/dev/urandom",
				toPath:   "out.txt",
				offset:   0,
				limit:    0,
			},
			wantErr: ErrUnsupportedFile,
		},
		{
			name: "when copy empty file",
			args: copyArgs{
				fromPath: "testdata/input_empty.txt",
				toPath:   "out.txt",
				offset:   0,
				limit:    0,
			},
			wantErr: ErrUnsupportedFile,
		},
	}
	for _, tt := range validateCases {
		t.Run("validate: " + tt.name, func(t *testing.T) {
			err := Copy(tt.args.fromPath, tt.args.toPath, tt.args.offset, tt.args.limit)
			assert.EqualError(t, err, tt.wantErr.Error(),)
		})
	}
}

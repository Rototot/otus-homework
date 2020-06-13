package main

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCopy(t *testing.T) {
	type args struct {
		fromPath string
		toPath   string
		offset   int64
		limit    int64
	}

	t.Run("simple", func(t *testing.T) {
		var args = *&args{
			fromPath: "testdata/input.txt",
			toPath:   "out.txt",
			offset:   0,
			limit:    0,
		}
		expectFile, err := os.Open(args.fromPath)
		assert.NoError(t, err)
		defer expectFile.Close()

		expectedFileInfo, err := expectFile.Stat()
		assert.NoError(t, err)

		result := Copy(args.fromPath, args.toPath, args.offset, args.limit)
		assert.NoError(t, result)
		assert.FileExists(t, args.toPath)

		actualFile, err := os.Open(args.toPath)
		assert.NoError(t, err)

		actualFileInfo, err := actualFile.Stat()
		assert.NoError(t, err)

		assert.Equal(t, expectedFileInfo.Size(), actualFileInfo.Size())

		if err := os.Remove(args.toPath); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("when set limit", func(t *testing.T) {
		var args = *&args{
			fromPath: "testdata/input.txt",
			toPath:   "out_limit_10.txt",
			offset:   0,
			limit:    10,
		}
		expectFile, err := os.Open("testdata/out_offset0_limit10.txt")
		assert.NoError(t, err)
		defer expectFile.Close()

		expectedFileInfo, err := expectFile.Stat()
		assert.NoError(t, err)

		result := Copy(args.fromPath, args.toPath, args.offset, args.limit)
		assert.NoError(t, result)
		assert.FileExists(t, args.toPath)

		actualFile, err := os.Open(args.toPath)
		assert.NoError(t, err)

		actualFileInfo, err := actualFile.Stat()
		assert.NoError(t, err)

		assert.Equal(t, expectedFileInfo.Size(), actualFileInfo.Size())

		if err := os.Remove(args.toPath); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("when set limit and offset", func(t *testing.T) {
		var args = *&args{
			fromPath: "testdata/input.txt",
			toPath:   "out_limit_10.txt",
			offset:   100,
			limit:    1000,
		}
		expectFile, err := os.Open("testdata/out_offset100_limit1000.txt")
		assert.NoError(t, err)
		defer expectFile.Close()

		expectedFileInfo, err := expectFile.Stat()
		assert.NoError(t, err)

		result := Copy(args.fromPath, args.toPath, args.offset, args.limit)
		assert.NoError(t, result)
		assert.FileExists(t, args.toPath)

		actualFile, err := os.Open(args.toPath)
		assert.NoError(t, err)

		actualFileInfo, err := actualFile.Stat()
		assert.NoError(t, err)

		assert.Equal(t, expectedFileInfo.Size(), actualFileInfo.Size())

		if err := os.Remove(args.toPath); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("when copy empty file", func(t *testing.T) {
		var args = *&args{
			fromPath: "testdata/input_empty.txt",
			toPath:   "out_limit_10.txt",
			offset:   0,
			limit:    1000,
		}
		expectFile, err := os.Open("testdata/out_empty.txt")
		assert.NoError(t, err)
		defer expectFile.Close()

		expectedFileInfo, err := expectFile.Stat()
		assert.NoError(t, err)

		result := Copy(args.fromPath, args.toPath, args.offset, args.limit)
		assert.NoError(t, result)
		assert.FileExists(t, args.toPath)

		actualFile, err := os.Open(args.toPath)
		assert.NoError(t, err)

		actualFileInfo, err := actualFile.Stat()
		assert.NoError(t, err)

		assert.Equal(t, expectedFileInfo.Size(), actualFileInfo.Size())

		if err := os.Remove(args.toPath); err != nil {
			t.Fatal(err)
		}
	})


	validateCases := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "limit < 0", args: args{
			fromPath: "testdata/input.txt",
			toPath:   "out.txt",
			offset:   -2,
			limit:    0,
		},
			wantErr: ErrInvalidArgument,
		},
		{
			name: "offset < 0", args: args{
			fromPath: "testdata/input.txt",
			toPath:   "out.txt",
			offset:   0,
			limit:    -1,
		},
			wantErr: ErrInvalidArgument,
		},
		{
			name: "empty source path", args: args{
			fromPath: "",
			toPath:   "out.txt",
			offset:   0,
			limit:    0,
		},
			wantErr: ErrEmptyPath,
		},
		{
			name: "empty dst path", args: args{
			fromPath: "testdata/input.txt",
			toPath:   "",
			offset:   0,
			limit:    0,
		},
			wantErr: ErrEmptyPath,
		},
		{
			name: "offset more than source size", args: args{
			fromPath: "testdata/input.txt",
			toPath:   "out.txt",
			offset:   1024 * 1024,
			limit:    0,
		},
			wantErr: ErrOffsetExceedsFileSize,
		},
		{
			name: "source path is directory", args: args{
			fromPath: "testdata",
			toPath:   "out.txt",
			offset:   0,
			limit:    0,
		},
			wantErr: errors.New("testdata is not a file"),
		},
	}
	for _, tt := range validateCases {
		t.Run(tt.name, func(t *testing.T) {
			err := Copy(tt.args.fromPath, tt.args.toPath, tt.args.offset, tt.args.limit)
			assert.Error(t, err)
			assert.EqualError(t, err, tt.wantErr.Error())
		})
	}

}

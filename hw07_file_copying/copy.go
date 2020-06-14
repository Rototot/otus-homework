package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrEmptyPath             = errors.New("empty file path")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrInvalidArgument       = os.ErrInvalid
)

type validator func() error

type copyArgs struct {
	fromPath string
	toPath   string
	offset   int64
	limit    int64
}

type fileOperation func(sourceFile *os.File, args copyArgs) error

func Copy(fromPath string, toPath string, offset, limit int64) error {
	var args = &copyArgs{
		fromPath: fromPath,
		toPath:   toPath,
		offset:   offset,
		limit:    limit,
	}

	// validate
	if err := validateArgs(*args); err != nil {
		return err
	}

	// open
	var sourceFile *os.File
	defer sourceFile.Close()

	var operations = []fileOperation{
		func(s *os.File, args copyArgs) error {
			file, err := os.Open(fromPath)
			sourceFile = file
			return err
		},
		seek,
		copying,
	}

	for _, action := range operations {
		if err := action(sourceFile, *args); err != nil {
			return err
		}
	}

	return nil
}

func seek(sourceFile *os.File, args copyArgs) error {
	fileInfo, err := sourceFile.Stat()
	if err != nil {
		return err
	}

	if args.offset > 0 {
		if args.offset > fileInfo.Size() {
			return ErrOffsetExceedsFileSize
		}
		// seek if need
		if _, err := sourceFile.Seek(args.offset, io.SeekStart); err != nil {
			return err
		}
	}

	return nil
}

func copying(sourceFile *os.File, args copyArgs) error {
	sourceFileInfo, err := sourceFile.Stat()
	if err != nil {
		return err
	}
	// create temp file
	tempDstFile, err := ioutil.TempFile("./", "out_*")
	if err != nil {
		return err
	}

	// calc max qty bytes for transfer
	var transferChunkSize int64 = 10 // 10 byte
	var transferLimit = sourceFileInfo.Size()
	if args.limit > 0 && args.limit < sourceFileInfo.Size() {
		transferLimit = args.limit
	}

	// bar init
	var refreshRate = 5 * time.Millisecond
	var bar = pb.Full.Start64(transferLimit)
	bar.Set(pb.Bytes, true)
	bar.Set(pb.Color, true)
	bar.SetRefreshRate(refreshRate)

	// copy
	var readBytes int64 = 0
	var chunkSize = transferChunkSize
	for readBytes < transferLimit {
		remainingBytes := transferLimit - readBytes
		if remainingBytes < chunkSize {
			chunkSize = remainingBytes
		}

		qty, err := io.CopyN(tempDstFile, sourceFile, chunkSize)

		readBytes += qty
		bar.Add64(qty)

		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		time.Sleep(refreshRate)
	}
	bar.Finish()

	if err := tempDstFile.Close(); err != nil {
		return err
	}

	// finish copy
	if err := os.Rename(tempDstFile.Name(), args.toPath); err != nil {
		return err
	}

	return nil
}

func validateArgs(args copyArgs) error {
	validators := []validator{
		func() error { return isNotEmptyPath(args.fromPath) },
		func() error { return isNotEmptyPath(args.toPath) },
		func() error { return isPositiveNumber(args.offset) },
		func() error { return isPositiveNumber(args.limit) },
		// is dir
		func() error {
			sourceFile, err := os.Open(args.fromPath)
			if err != nil {
				return err
			}
			defer sourceFile.Close()

			fileInfo, err := sourceFile.Stat()
			if err != nil {
				return err
			}

			if fileInfo.IsDir() {
				return fmt.Errorf("%s is not a file", args.fromPath)
			}

			return nil
		},
	}
	if err := validate(validators); err != nil {
		return err
	}

	return nil
}

func validate(validators []validator) error {
	var err error
	for _, v := range validators {
		err = v()
		if err != nil {
			break
		}
	}

	return err
}

func isNotEmptyPath(path string) error {
	if path == "" {
		return ErrEmptyPath
	}

	return nil
}

func isPositiveNumber(value int64) error {
	if value < 0 {
		return ErrInvalidArgument
	}

	return nil
}

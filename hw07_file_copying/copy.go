package main

import (
	"errors"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	copyFile, err := os.OpenFile(fromPath, os.O_RDONLY, 0777)
	defer copyFile.Close()
	size, err := copyFile.Seek(0, io.SeekEnd)
	if err != nil {
		return ErrUnsupportedFile
	}
	if offset > size {
		return ErrOffsetExceedsFileSize
	}
	_, err = copyFile.Seek(offset, io.SeekStart)
	targetFile, err := os.Create(toPath)
	defer targetFile.Close()
	if limit == 0 {
		_, err = io.Copy(targetFile, copyFile)
	} else {
		_, err = io.CopyN(targetFile, copyFile, limit)
	}
	return nil
}

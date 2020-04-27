package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrInputParams           = errors.New("from and to are required")
)

func Copy(fromPath string, toPath string, offset, limit int64) error {
	if fromPath == "" || toPath == "" {
		fmt.Println("params:", fromPath, toPath)
		return ErrInputParams
	}
	fromPathSize, err := getFileSize(fromPath)
	if err != nil {
		return fmt.Errorf("getting file size error: %s\n", err)
	}
	if offset > fromPathSize {
		return ErrOffsetExceedsFileSize
	}

	var n int64
	if limit > 0 {
		n = limit
	} else {
		n = fromPathSize
	}

	fileFrom, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("open \"from\" file error: %s\n", err)
	}
	defer fileFrom.Close()

	if offset > 0 {
		fileFrom.Seek(offset, io.SeekStart)
	} else if offset < 0 {
		fileFrom.Seek(offset, io.SeekEnd)
	}

	fileTo, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("create \"to\" file error: %s\n", err)
	}
	defer fileTo.Close()

	pb := NewProxyReader(fileFrom, n)
	_, err = io.CopyN(fileTo, pb, n)
	pb.Finish()

	if err != nil && err != io.EOF {
		return fmt.Errorf("copy file error: %s\n", err)
	}

	return nil
}

func getFileSize(filepath string) (int64, error) {
	fi, err := os.Stat(filepath)
	if err != nil {
		return 0, err
	}
	size := fi.Size()
	if size == 0 {
		return 0, ErrUnsupportedFile
	}
	return size, nil
}

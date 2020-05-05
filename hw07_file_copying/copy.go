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
	ErrInputParams           = errors.New("from and to params are required")
)

func Copy(fromPath string, toPath string, offset, limit int64) error {
	if fromPath == "" || toPath == "" {
		return ErrInputParams
	}
	inFileSize, err := getFileSize(fromPath)
	if err != nil {
		return err
	}
	if offset > inFileSize {
		return ErrOffsetExceedsFileSize
	}

	inFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("open \"in\" file error: %s\n", err)
	}
	defer inFile.Close()

	if offset > 0 {
		inFile.Seek(offset, io.SeekStart)
	} else if offset < 0 {
		inFile.Seek(offset, io.SeekEnd)
	}

	outFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("create \"out\" file error: %s\n", err)
	}
	defer outFile.Close()

	// если бы не нужно было реализовывать прогрессбар, то можно просто взять
	// минимальное значение из (inFileSize и limit) при условии что limit > 0
	// а в функции CopyN не реагировать на io.EOF
	n := getBytesCount(inFileSize, offset, limit)
	pb := NewProxyReader(inFile, n)
	fmt.Printf("\nCopying from '%s' to '%s'\n", fromPath, toPath)
	_, err = io.CopyN(outFile, pb, n)
	pb.Finish()

	if err != nil {
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

func getBytesCount(fileSize, offset, limit int64) int64 {
	start := offset
	if offset < 0 {
		start = fileSize + offset
	}
	end := fileSize
	endLimited := start + limit
	if limit > 0 && end > endLimited {
		end = endLimited
	}
	return (end - start)
}

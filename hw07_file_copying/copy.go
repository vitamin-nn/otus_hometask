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
	ErrSeekSetting           = errors.New("seek setting error")
	ErrOpenInFile            = errors.New("open \"in\" file error")
	ErrOutFileCreate         = errors.New("create \"out\" file error")
	ErrCopyFile              = errors.New("copy file error")
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
		return ErrOpenInFile
	}
	defer inFile.Close()

	seekMode := io.SeekStart
	if offset < 0 {
		seekMode = io.SeekEnd
	}
	if offset != 0 {
		if _, err = inFile.Seek(offset, seekMode); err != nil {
			return ErrSeekSetting
		}
	}

	outFile, err := os.Create(toPath)
	if err != nil {
		return ErrOutFileCreate
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
		return ErrCopyFile
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

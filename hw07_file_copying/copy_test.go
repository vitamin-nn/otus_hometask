package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var last99 = `The Go Gopher
Copyright Terms of Service Privacy Policy Report a website issue
Supported by Google
`

func TestCopy(t *testing.T) {
	t.Run("empty from", func(t *testing.T) {
		result := Copy("", "/dev/null", 0, 0)
		require.Equal(t, ErrInputParams, result)
	})

	t.Run("unsupported file", func(t *testing.T) {
		result := Copy("/dev/urandom", "/dev/null", 0, 0)
		require.Equal(t, ErrUnsupportedFile, result)
	})

	t.Run("incorrect offset", func(t *testing.T) {
		result := Copy("testdata/input.txt", "/dev/null", 10000, 0)
		require.Equal(t, ErrOffsetExceedsFileSize, result)
	})

	t.Run("negative offset", func(t *testing.T) {
		tmpfile, err := ioutil.TempFile("", "hw07_test_file_")
		if err != nil {
			t.Fatal(err)
		}
		tmpFileName := tmpfile.Name()
		defer os.Remove(tmpFileName)
		err = tmpfile.Close()
		if err != nil {
			t.Fatal(err)
		}

		outFile, err := os.Open(tmpFileName)
		if err != nil {
			t.Fatal(err)
		}

		defer outFile.Close()
		result := Copy("testdata/input.txt", tmpFileName, -99, 0)
		require.Nil(t, result)

		b, err := ioutil.ReadAll(outFile)
		if err != nil {
			t.Fatal(err)
		}

		require.Equal(t, last99, string(b))
	})
}

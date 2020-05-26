package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("cannot run cmd", func(t *testing.T) {
		cmd := []string{"/dev/null", "1", "2"}
		env := Environment{"TEST": "VALUE"}
		code := RunCmd(cmd, env)
		require.Equal(t, defaultExitErrCode, code)
	})

	t.Run("cmd error code return", func(t *testing.T) {
		script := `#!/bin/sh
exit 1
`
		tmpFileName, err := makeShFile(script)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpFileName)

		cmd := []string{
			"/bin/sh",
			tmpFileName,
		}
		env := Environment{}
		code := RunCmd(cmd, env)
		require.Equal(t, 1, code)
	})

	t.Run("cmd writes to Stdout", func(t *testing.T) {
		script := `#!/bin/sh
printf "out"
`
		tmpFileName, err := makeShFile(script)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpFileName)

		rescueStdout := os.Stdout
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}
		os.Stdout = w
		cmd := []string{
			"/bin/sh",
			tmpFileName,
		}
		env := Environment{}
		_ = RunCmd(cmd, env)
		w.Close()
		out, err := ioutil.ReadAll(r)
		if err != nil {
			t.Fatal(err)
		}
		os.Stdout = rescueStdout
		require.Equal(t, "out", string(out))
	})
}
func makeShFile(content string) (string, error) {
	tmpfile, err := ioutil.TempFile("", "hw08_test_file_*.sh")
	if err != nil {
		return "", err
	}
	tmpFileName := tmpfile.Name()

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		os.Remove(tmpFileName)
		return "", err
	}
	if err := tmpfile.Close(); err != nil {
		os.Remove(tmpFileName)
		return "", err
	}
	return tmpFileName, nil
}

package main

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const wrongNameSymbol = "="

type Environment map[string]string

var (
	ErrOpenDir = errors.New("open dir error")
)

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, ErrOpenDir
	}

	env := Environment{}
	for _, file := range files {
		if file.IsDir() || isWrongName(file.Name()) {
			continue
		}
		envVarName := file.Name()
		if file.Size() == 0 {
			os.Unsetenv(envVarName)
			continue
		}
		s, err := readLine(filepath.Join(dir, file.Name()))
		if err != nil {
			continue
		}
		s = filterStrValue(s)
		if _, ok := os.LookupEnv(envVarName); ok {
			os.Unsetenv(envVarName)
			os.Setenv(envVarName, s)
		} else {
			env[envVarName] = s
		}
	}
	return env, nil
}

func isWrongName(name string) bool {
	return strings.Contains(name, wrongNameSymbol)
}

func readLine(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()

	return scanner.Text(), scanner.Err()
}

func filterStrValue(s string) string {
	s = strings.TrimRight(s, " \t")
	s = strings.Replace(s, "\x00", "\n", -1)
	return s
}

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("incorrect dir", func(t *testing.T) {
		_, err := ReadDir("/dev/null")
		require.NotNil(t, err)
		require.Equal(t, ErrOpenDir, err)
	})
	// проверяем, что переменная именно ансетится,
	// а не заменяется значение пустой строкой
	t.Run("unset env", func(t *testing.T) {
		d, err := ioutil.TempDir("", "env_test")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(d)
		uns := "UNSET"
		err = makeFile(filepath.Join(d, uns), "")
		if err != nil {
			t.Fatal(err)
		}

		os.Setenv(uns, "test")
		_, err = ReadDir(d)
		require.Nil(t, err)
		_, ok := os.LookupEnv(uns)
		require.False(t, ok)
	})
	t.Run("incorect variable name", func(t *testing.T) {
		d, err := ioutil.TempDir("", "env_test")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(d)

		value := "test"
		ninc := "INC=ORRECT"
		err = makeFile(filepath.Join(d, ninc), value)
		if err != nil {
			t.Fatal(err)
		}

		ncorr := "CORRECT"
		err = makeFile(filepath.Join(d, ncorr), value)
		if err != nil {
			t.Fatal(err)
		}

		result, err := ReadDir(d)
		require.Nil(t, err)
		_, ok := os.LookupEnv(ninc)
		require.False(t, ok)

		v, ok := result[ncorr]
		require.True(t, ok)
		require.Equal(t, value, v)
	})
}

func makeFile(f string, val string) error {
	v := []byte(val)
	err := ioutil.WriteFile(f, v, 0644)
	if err != nil {
		return err
	}
	return nil
}

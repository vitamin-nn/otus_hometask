package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	d, err := ioutil.TempDir("", "parser_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(d)
	filename := filepath.Join(d, "parser_test_model.go")
	err = makeTestFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("check parsed data", func(t *testing.T) {
		parsedData, err := parse(filename)
		require.Nil(t, err)
		require.Equal(t, parsedData.PackageName, "parser")
		for structName, structFields := range parsedData.Structures {
			require.Equal(t, structName, "TestParserStruct")

			var i = 0
			for _, field := range structFields {
				if field.Name == "ID" {
					require.True(t, field.Len.IsSet)
					len, ok := field.Len.Value.(int)
					require.True(t, ok)
					require.Equal(t, len, 36)
					i++
				}
				if field.Name == "Age" {
					require.True(t, field.Min.IsSet)
					min, ok := field.Min.Value.(int)
					require.True(t, ok)
					require.Equal(t, min, 18)
					require.True(t, field.Max.IsSet)
					max, ok := field.Max.Value.(int)
					require.True(t, ok)
					require.Equal(t, max, 50)
					i++
				}
				if field.Name == "Email" {
					require.True(t, field.Regexp.IsSet)
					regexp, ok := field.Regexp.Value.(string)
					require.True(t, ok)
					require.Equal(t, regexp, "^.+$")
					i++
				}
				if field.Name == "Role" {
					require.True(t, field.Enum.IsSet)
					enum := field.GetEnumList()
					require.NotNil(t, enum)
					eq := reflect.DeepEqual(enum, []string{"\"admin\"", "\"stuff\""})
					require.True(t, eq)
					require.Equal(t, field.Type.Name, "string")
					i++
				}
				if field.Name == "Phones" {
					require.True(t, field.Len.IsSet)
					require.True(t, field.Type.IsSlice)
					require.Equal(t, field.Type.Name, "string")
					i++
				}
			}
			require.Equal(t, i, 5)
		}

	})
}

func makeTestFile(f string) error {
	v := []byte(fileContent)
	err := ioutil.WriteFile(f, v, 0644)
	if err != nil {
		return err
	}
	return nil
}

var fileContent = `package parser

type UserRole string

// NOTE: Several struct specs in one type declaration are allowed

type TestParserStruct struct {
	ID     string 	` + "`json:\"id\" validate:\"len:36\"`" + `
	Name   string
	Age    int      ` + "`validate:\"min:18|max:50\"`" + `
	Email  string   ` + "`validate:\"regexp:^.+$\"`" + `
	Role   UserRole ` + "`validate:\"in:admin,stuff\"`" + `
	Phones []string ` + "`validate:\"len:11\"`" + `
}

type Token struct {
	Header    []byte ` + "`json:\"header\"`" + `
}
`

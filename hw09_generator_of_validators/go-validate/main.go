package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

const (
	minArgs          = 2
	parseTagName     = "validate"
	generatedPostfix = "_validation_generated.go"
)

var (
	ErrInvalidArgument = errors.New("invalid argument")
)

func main() {
	args := os.Args
	if len(args) != minArgs {
		log.Fatalln(ErrInvalidArgument)
	}

	inFilePath := args[1]
	parsedData, err := parse(inFilePath)
	if err != nil {
		log.Fatalln(err)
	}
	outFilePath := getGeneratedFilePath(inFilePath)
	err = generate(outFilePath, parsedData)
	if err != nil {
		log.Fatalln(err)
	}
}

func getGeneratedFilePath(inputFilePath string) string {
	dir := filepath.Dir(inputFilePath)
	filename := filepath.Base(inputFilePath)
	re := regexp.MustCompile(`(.+)\..*$`)
	chanks := re.FindStringSubmatch(filename)
	generatedName := filename + generatedPostfix
	if n := chanks[1]; n != "" {
		generatedName = n + generatedPostfix
	}
	return filepath.Join(dir, generatedName)
}

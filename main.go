package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/moby/buildkit/frontend/dockerfile/dockerignore"
	"github.com/moby/moby/pkg/fileutils"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Error: dockerignore path missing")
		printUsage()
		os.Exit(1)
	}

	if os.Args[1] == "--help" || os.Args[1] == "-h" {
		printUsage()
		return
	}

	ignoreFile, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer ignoreFile.Close()

	if err := process(ignoreFile, os.Stdin, os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func process(ignoreReader io.Reader, inputPathreader io.Reader, outputPathWriter io.Writer) error {
	ignorePatterns, err := dockerignore.ReadAll(ignoreReader)
	if err != nil {
		return err
	}

	ignorePatternMatcher, err := fileutils.NewPatternMatcher(ignorePatterns)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(inputPathreader)
	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	for scanner.Scan() {
		path := scanner.Text()

		normalizedPath := filepath.Clean(path)
		if filepath.IsAbs(normalizedPath) {
			normalizedPath, _ = filepath.Rel(workingDir, normalizedPath)
		}

		matches, err := ignorePatternMatcher.Matches(normalizedPath)
		if err != nil {
			return err
		}

		if !matches {
			if _, err := fmt.Fprintln(outputPathWriter, path); err != nil {
				return err
			}
		}
	}

	return scanner.Err()
}

func printUsage() {
	fmt.Fprintln(os.Stderr, `
Filters lines of paths in stdin based on a .dockerignore file.

Usage: dockerignore-filter <dockerignore-path>
	`)
}

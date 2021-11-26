package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

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

	ignorePatterns, err := loadIgnoreFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	ignorePatternMatcher, err := fileutils.NewPatternMatcher(ignorePatterns)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		matches, err := ignorePatternMatcher.Matches(text)
		if err != nil {
			log.Fatal(err)
		}

		if !matches {
			if _, err := fmt.Println(text); err != nil {
				log.Fatal(err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func loadIgnoreFile(filePath string) ([]string, error) {
	fileHandle, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer fileHandle.Close()
	return dockerignore.ReadAll(fileHandle)
}

func printUsage() {
	fmt.Fprintln(os.Stderr, `
Filters lines of paths in stdin based on a .dockerignore file.

Usage: dockerignore-filter <dockerignore-path>
	`)
}

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const webpackPrefix = "webpack://"

type mappingStruct struct {
	Version        int      `json:"version"`
	Sources        []string `json:"sources"`
	Names          []string `json:"names"`
	Mappings       string   `json:"mappings"`
	File           string   `json:"file"`
	SourcesContent []string `json:"sourcesContent"`
	SourceRoot     string   `json:"sourceRoot"`
}

func main() {
	filename := flag.String("f", "", "Input filename, or use stdin if not present")
	summary := flag.Bool("s", false, "Output summary")
	outputDir := flag.String("dir", ".", "Output extracted source to this directory")
	index := flag.Int("i", -1, "Array index to output")
	showName := flag.Bool("n", false, "Output source name")
	showContent := flag.Bool("c", false, "Output sourcesContent")
	showSources := flag.Bool("sources", false, "Show all source names")
	doOutputFiles := flag.Bool("o", false, "Output one file, if index has been set, or all files, if it hasn't")
	flag.Parse()

	var inputFile *os.File
	if *filename == "" {
		inputFile = os.Stdin
	} else {
		var err error
		inputFile, err = os.Open(*filename)
		if err != nil {
			panic(err)
		}
	}

	var contents []byte

	contents, err := ioutil.ReadAll(inputFile)
	if err != nil {
		panic(err)
	}

	var mapping mappingStruct
	if err := json.Unmarshal(contents, &mapping); err != nil {
		panic(err)
	}

	if *summary {
		if *filename == "" {
			fmt.Printf("(stdio)\n")
		} else {
			fmt.Printf("%v\n", *filename)
		}

		fmt.Printf("Version:        %v\n", mapping.Version)
		fmt.Printf("Sources:        %v\n", len(mapping.Sources))
		fmt.Printf("Names:          %v\n", len(mapping.Names))
		fmt.Printf("Mappings:       %v characters\n", len(mapping.Mappings))
		fmt.Printf("File:           %q\n", mapping.File)
		fmt.Printf("SourcesContent: %v\n", len(mapping.SourcesContent))
		fmt.Printf("SourceRoot:     %q\n", mapping.SourceRoot)
	}

	if *showSources {
		for i, name := range mapping.Sources {
			fmt.Printf("[%v] = %q\n", i, name)
		}
	}

	if *showName {
		fmt.Printf("%q\n", mapping.Sources[*index])
	}

	if *showContent {
		fmt.Printf("%s\n", mapping.SourcesContent[*index])
	}

	if *doOutputFiles {
		// Output either all files or just one
		if *index >= 0 {
			outputFiles(mapping, *index, *outputDir)
		} else {
			for index := range mapping.Sources {
				outputFiles(mapping, index, *outputDir)
			}
		}
	}
}

func removeWebpackPrefix(filename string) string {
	if strings.Index(filename, webpackPrefix) == 0 {
		return filename[len(webpackPrefix):]
	} else {
		return filename
	}
}

func outputFiles(mapping mappingStruct, index int, outputDir string) {
	filename := removeWebpackPrefix(mapping.Sources[index])
	fullPath := filepath.Join(outputDir, filename)

	directory := filepath.Dir(fullPath)
	if err := os.MkdirAll(directory, os.FileMode(0755)); err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(fullPath, []byte(mapping.SourcesContent[index]), 0644); err != nil {
		panic(err)
	}
}

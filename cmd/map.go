/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
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

var filename *string
var summary *bool
var outputDir *string
var index *int
var showName *bool
var showContent *bool
var showSources *bool
var doOutputFiles *bool

// mapCmd represents the map command
var mapCmd = &cobra.Command{
	Use: "map",
	// 	Short: "A brief description of your command",
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

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

	},
}

func init() {
	rootCmd.AddCommand(mapCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mapCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mapCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	filename = mapCmd.Flags().StringP("file", "f", "", "Input filename, or use stdin if not present")
	summary = mapCmd.Flags().BoolP("summary", "s", false, "Output summary")
	outputDir = mapCmd.Flags().String("dir", ".", "Output extracted source to this directory")
	index = mapCmd.Flags().IntP("index", "i", -1, "Array index to output")
	showName = mapCmd.Flags().BoolP("name", "n", false, "Output source name")
	showContent = mapCmd.Flags().BoolP("content", "c", false, "Output sourcesContent")
	showSources = mapCmd.Flags().Bool("sources", false, "Show all source names")
	doOutputFiles = mapCmd.Flags().BoolP("onefile", "o", false, "Output one file, if index has been set, or all files, if it hasn't")
}

func removeWebpackPrefix(filename string) string {
	if strings.Index(filename, webpackPrefix) == 0 {
		return filename[len(webpackPrefix):]
	}
	return filename
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

package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/onethousandone/todos/parser"
)

var (
	pkgs       = flag.String("pkgs", "", "A list of comma separated package paths")
	outputFile = flag.String("output-file", "TODO.md", "Output file for todo's")
)

func main() {
	var absDirs []string
	var todos []parser.Todo
	flag.Parse()

	if len(*pkgs) == 0 {
		log.Fatalf("the flag -pkgs must be set")
	}
	packageDirs := strings.Split(*pkgs, ",")

	// Look for absolute paths.
	for _, dir := range packageDirs {
		dir, err := filepath.Abs(dir)
		if err != nil {
			log.Fatalf("unable to determine absolute filepath for requested path %s: %v", dir, err)
		}
		absDirs = append(absDirs, dir)
	}

	// Parse every directory containing a go-package.
	for _, dir := range absDirs {
		prg, err := parser.ParsePackage(dir)
		if err != nil {
			log.Fatalf("parsing program: %v", err)
		}
		todos = append(todos, parser.GetTodos(prg)...)
	}

	// Write analysis to markdown file.
	var analysis = struct {
		Command string
		Todos   []parser.Todo
	}{
		Command: strings.Join(os.Args[1:], " "),
		Todos:   todos,
	}

	var buf bytes.Buffer
	if err := generatedTmpl.Execute(&buf, analysis); err != nil {
		log.Fatalf("generating code: %v", err)
	}

	outputPath := *outputFile
	if err := ioutil.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		log.Fatalf("writing output: %s", err)
	}
}

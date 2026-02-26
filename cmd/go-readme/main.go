// Command go-readme generates a README.md for a Go package.
//
// Usage:
//
//	go-readme [flags]
//
// Flags:
//
//	-dir string
//	    directory of the Go package to document (default ".")
//	-output string
//	    output file path (default "README.md")
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sermachage/go-readme/internal/analyzer"
	"github.com/sermachage/go-readme/internal/generator"
)

func main() {
	dir := flag.String("dir", ".", "directory of the Go package to document")
	output := flag.String("output", "README.md", "output file path")
	flag.Parse()

	pkg, err := analyzer.Analyze(*dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "go-readme: analyze: %v\n", err)
		os.Exit(1)
	}

	content, err := generator.Generate(pkg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "go-readme: generate: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(*output, []byte(content), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "go-readme: write %s: %v\n", *output, err)
		os.Exit(1)
	}

	fmt.Printf("go-readme: wrote %s\n", *output)
}

package main

import (
	"fmt"
	"grep/pkg/config"
	"grep/pkg/grep"
	"grep/pkg/parser"
	"grep/pkg/reader"
	"os"
)

func main() {
	cfg := &config.Config{}

	// Parse command-line arguments
	err := parser.ParseArgs(cfg)
	if err != nil {
		os.Exit(1)
	}

	// Read input lines
	lines, err := reader.ReadLines(cfg.FilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	// Perform grep operation
	results, err := grep.Grep(lines, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during grep: %v\n", err)
		os.Exit(1)
	}

	// Format and output results
	output := grep.FormatResults(results, cfg)
	for _, line := range output {
		fmt.Println(line)
	}
}

package parser

import (
	"fmt"
	"grep/pkg/config"
	"os"

	"github.com/akamensky/argparse"
)

func ParseArgs(config *config.Config) error {
	parser := argparse.NewParser("grep", "Utility for filtering text stream (analog of grep command)")

	ContextAfter := parser.Int("A", "after-context", &argparse.Options{Default: 0, Help: "Print N lines of trailing context after matching lines", Required: false})
	ContextBefore := parser.Int("B", "before-context", &argparse.Options{Default: 0, Help: "Print N lines of leading context before matching lines", Required: false})
	ContextStrings := parser.Int("C", "context", &argparse.Options{Default: 0, Help: "Print N lines of output context", Required: false})
	IsCountMatching := parser.Flag("c", "count", &argparse.Options{Default: false, Help: "Only print a count of matching lines", Required: false})
	IsIgnoreRegister := parser.Flag("i", "ignore-case", &argparse.Options{Default: false, Help: "Ignore case distinctions", Required: false})
	IsInvertOutput := parser.Flag("v", "invert-match", &argparse.Options{Default: false, Help: "Invert the sense of matching, to select non-matching lines", Required: false})
	IsFixed := parser.Flag("F", "fixed-strings", &argparse.Options{Default: false, Help: "Interpret pattern as a fixed string, not a regular expression", Required: false})
	IsNumerableStrings := parser.Flag("n", "line-number", &argparse.Options{Default: false, Help: "Prefix each line of output with the line number", Required: false})

	Pattern := parser.String("", "pattern", &argparse.Options{Help: "Pattern to search for", Required: true})
	FilePath := parser.String("", "file", &argparse.Options{Help: "File to search in (if not provided, reads from STDIN)", Required: false, Default: ""})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return err
	}

	config.ContextAfter = *ContextAfter
	config.ContextBefore = *ContextBefore
	config.ContextStrings = *ContextStrings
	config.IsCountMatching = *IsCountMatching
	config.IsIgnoreRegister = *IsIgnoreRegister
	config.IsInvertOutput = *IsInvertOutput
	config.IsFixed = *IsFixed
	config.IsNumerableStrings = *IsNumerableStrings
	config.Pattern = *Pattern
	config.FilePath = *FilePath

	if config.ContextStrings > 0 {
		config.ContextAfter = config.ContextStrings
		config.ContextBefore = config.ContextStrings
	}

	return nil
}

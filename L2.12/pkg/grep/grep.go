package grep

import (
	"fmt"
	"grep/pkg/config"
	"regexp"
	"strings"
)

type Result struct {
	LineNumber int
	Line       string
	IsMatch    bool
}

func Grep(lines []string, cfg *config.Config) ([]Result, error) {
	var matcher func(string) bool
	var err error

	if cfg.IsFixed {
		pattern := cfg.Pattern
		if cfg.IsIgnoreRegister {
			pattern = strings.ToLower(pattern)
			matcher = func(line string) bool {
				return strings.Contains(strings.ToLower(line), pattern)
			}
		} else {
			matcher = func(line string) bool {
				return strings.Contains(line, pattern)
			}
		}
	} else {
		pattern := cfg.Pattern
		if cfg.IsIgnoreRegister {
			pattern = "(?i)" + pattern
		}
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid regular expression: %w", err)
		}
		matcher = func(line string) bool {
			return re.MatchString(line)
		}
	}

	matchingLines := make(map[int]bool)
	for i, line := range lines {
		matches := matcher(line)
		if cfg.IsInvertOutput {
			matches = !matches
		}
		if matches {
			matchingLines[i] = true
		}
	}

	if cfg.IsCountMatching {
		return []Result{{Line: fmt.Sprintf("%d", len(matchingLines))}}, nil
	}

	results := buildResultsWithContext(lines, matchingLines, cfg)

	return results, err
}

func buildResultsWithContext(lines []string, matchingLines map[int]bool, cfg *config.Config) []Result {
	linesToInclude := make(map[int]bool)
	contextLines := make(map[int]bool)

	for lineNum := range matchingLines {
		linesToInclude[lineNum] = true

		for i := 1; i <= cfg.ContextBefore; i++ {
			beforeLine := lineNum - i
			if beforeLine >= 0 {
				if !matchingLines[beforeLine] {
					contextLines[beforeLine] = true
				}
				linesToInclude[beforeLine] = true
			}
		}

		for i := 1; i <= cfg.ContextAfter; i++ {
			afterLine := lineNum + i
			if afterLine < len(lines) {
				if !matchingLines[afterLine] {
					contextLines[afterLine] = true
				}
				linesToInclude[afterLine] = true
			}
		}
	}

	var results []Result
	for i := 0; i < len(lines); i++ {
		if linesToInclude[i] {
			result := Result{
				LineNumber: i + 1,
				Line:       lines[i],
				IsMatch:    matchingLines[i],
			}
			results = append(results, result)
		}
	}

	return results
}

func FormatResults(results []Result, cfg *config.Config) []string {
	var output []string

	for _, result := range results {
		var line string

		if cfg.IsNumerableStrings {
			line = fmt.Sprintf("%d:%s", result.LineNumber, result.Line)
		} else {
			line = result.Line
		}

		output = append(output, line)
	}

	return output
}

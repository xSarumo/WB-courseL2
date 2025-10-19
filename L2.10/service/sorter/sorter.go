package sorter

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"sort/external/model"
	"strconv"
	"strings"
)

var monthMap = map[string]int{
	"jan": 1, "feb": 2, "mar": 3, "apr": 4,
	"may": 5, "jun": 6, "jul": 7, "aug": 8,
	"sep": 9, "oct": 10, "nov": 11, "dec": 12,
}

func SortStrings(config *model.Config) {
	lines, err := readLines(config.IncludePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in read file: %v", err)
		os.Exit(1)
	}

	less := func(i, j int) bool {
		valI := getColumnValue(lines[i], config.ColumnSort)
		valJ := getColumnValue(lines[j], config.ColumnSort)

		if config.Numeric {
			numI, errI := strconv.ParseFloat(valI, 64)
			numJ, errJ := strconv.ParseFloat(valJ, 64)
			if errI == nil && errJ == nil {
				return numI < numJ
			}
		}
		if config.MonthSort {
			monthI, okI := monthMap[strings.ToLower(valI)]
			monthJ, okJ := monthMap[strings.ToLower(valJ)]

			if okI && okJ {
				return monthI < monthJ
			}
		}
		return valI < valJ
	}

	sort.Slice(lines, less)

	if config.ReverseSort {
		reversSort(lines)
	}

	if config.UniqueStrings {
		lines = unique(lines)
	}

	printLines(lines)
}

func readLines(filePath string) ([]string, error) {
	var reader io.Reader
	if filePath == "" {
		reader = os.Stdin
	} else {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		reader = file
	}

	var lines []string
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func getColumnValue(line string, col int) string {
	if col <= 1 {
		return line
	}
	fields := strings.Split(line, "\t")
	if col-1 < len(fields) {
		return fields[col-1]
	}
	return ""
}

func unique(lines []string) []string {
	if len(lines) == 0 {
		return lines
	}
	result := []string{lines[0]}
	for i := 1; i < len(lines); i++ {
		if lines[i] != lines[i-1] {
			result = append(result, lines[i])
		}
	}
	return result
}

func printLines(lines []string) {
	for _, line := range lines {
		fmt.Println(line)
	}
}

func reversSort(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

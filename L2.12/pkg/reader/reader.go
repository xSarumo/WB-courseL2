package reader

import (
	"bufio"
	"io"
	"os"
)

func ReadLines(filePath string) ([]string, error) {
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

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

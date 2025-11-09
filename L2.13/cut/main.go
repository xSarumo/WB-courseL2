package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

type fieldsSpec = []int

func parseFields(spec string) (fieldsSpec, error) {
	if strings.TrimSpace(spec) == "" {
		return nil, errors.New("empty fields specification")
	}

	seen := make(map[int]struct{})
	for _, part := range strings.Split(spec, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if i := strings.IndexByte(part, '-'); i >= 0 {
			a := strings.TrimSpace(part[:i])
			b := strings.TrimSpace(part[i+1:])
			if a == "" || b == "" {
				return nil, fmt.Errorf("invalid range: %q", part)
			}
			start, err := strconv.Atoi(a)
			if err != nil || start <= 0 {
				return nil, fmt.Errorf("invalid range start: %q", a)
			}
			end, err := strconv.Atoi(b)
			if err != nil || end <= 0 {
				return nil, fmt.Errorf("invalid range end: %q", b)
			}
			if start > end {
				return nil, fmt.Errorf("invalid range %d-%d: start > end", start, end)
			}
			for k := start; k <= end; k++ {
				seen[k] = struct{}{}
			}
		} else {
			n, err := strconv.Atoi(part)
			if err != nil || n <= 0 {
				return nil, fmt.Errorf("invalid field number: %q", part)
			}
			seen[n] = struct{}{}
		}
	}

	if len(seen) == 0 {
		return nil, errors.New("no fields parsed")
	}

	out := make([]int, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	sort.Ints(out)
	return out, nil
}

func processLine(line string, fields fieldsSpec, delim string, separatedOnly bool) (string, bool) {
	if separatedOnly && !strings.Contains(line, delim) {
		return "", false
	}

	parts := strings.Split(line, delim)
	var out []string
	for _, f := range fields {
		idx := f - 1
		if idx >= 0 && idx < len(parts) {
			out = append(out, parts[idx])
		}
	}
	if len(out) == 0 {
		return "", false
	}
	return strings.Join(out, delim), true
}

func readAndProcess(r io.Reader, w io.Writer, fields fieldsSpec, delim string, separatedOnly bool) error {
	scanner := bufio.NewScanner(r)
	const maxCapacity = 10 * 1024 * 1024
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, maxCapacity)

	outWriter := bufio.NewWriter(w)
	defer outWriter.Flush()

	for scanner.Scan() {
		line := scanner.Text()
		if s, ok := processLine(line, fields, delim, separatedOnly); ok {
			_, _ = outWriter.WriteString(s)
			_, _ = outWriter.WriteString("\n")
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func main() {
	fSpec := flag.String("f", "", "fields to print (e.g. '1,3-5') (required)")
	d := flag.String("d", "\t", "delimiter (default tab)")
	s := flag.Bool("s", false, "only lines with delimiter are printed")
	help := flag.Bool("help", false, "show help")
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	if strings.TrimSpace(*fSpec) == "" {
		fmt.Fprintln(os.Stderr, "Error: -f fields is required")
		flag.Usage()
		os.Exit(1)
	}

	fields, err := parseFields(*fSpec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing -f: %v\n", err)
		os.Exit(2)
	}

	delim := *d
	if delim == "\\t" {
		delim = "\t"
	}
	if delim == "" {
		fmt.Fprintln(os.Stderr, "Error: delimiter cannot be empty")
		os.Exit(2)
	}

	if err := readAndProcess(os.Stdin, os.Stdout, fields, delim, *s); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

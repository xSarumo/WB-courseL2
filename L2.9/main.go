package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func UnpackString(str string) (string, error) {
	if str == "" {
		return "", nil
	}

	var strBuild strings.Builder
	length := len(str)
	escapeFlag := false
	for i := 0; i < length; i++ {
		if rune(str[i]) == '\\' {
			escapeFlag = true
			continue
		}

		if unicode.IsLetter(rune(str[i])) || escapeFlag {
			repeat := 1
			if i+1 < length && unicode.IsDigit(rune(str[i+1])) {
				repeat, _ = strconv.Atoi(string(str[i+1]))
			}
			strBuild.WriteString(strings.Repeat(string(rune(str[i])), repeat))
			escapeFlag = false
		}
	}

	result := strBuild.String()

	if result == "" {
		return "", fmt.Errorf("Uncorrect string format!!!")
	}

	return result, nil
}

func main() {
	fmt.Println("Hello from main.go")
}
